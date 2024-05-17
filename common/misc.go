package common

import (
	"fmt"
	. "github.com/knaka/go-utils"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/tidwall/gjson"
	"net/url"
	"os"
)

var DirsToCleanUp []string

// Clean cleans up generated files.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Clean() {
	fmt.Println("Cleaning...")
	for _, dir := range DirsToCleanUp {
		V0(os.RemoveAll(dir))
	}
}

//goland:noinspection GoUnusedExportedType, GoUnnecessarilyExportedIdentifiers
type Env mg.Namespace

// Compose prints text in the .env format that references the Docker Compose configuration.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Env) Compose() {
	json := V(sh.Output("docker", "compose", "config", "--format", "json"))
	//s, _ := exec.LookPath("docker-compose")
	//fmt.Printf("docker-compose: %s\n", s)
	//dockerId := V(sh.Output("/Applications/Docker.app/Contents/Resources/bin/docker-compose", "ps", "--quiet", "db"))
	//cmd := exec.Command("sh", "-c", "/Applications/Docker.app/Contents/Resources/bin/docker compose ps --quiet db")
	//out, err := cmd.Output()
	//if err != nil {
	//	fmt.Printf("Error: %s\n", err)
	//	fmt.Printf("Output: %s\n", out)
	//	return
	//}
	//// 出力結果を表示
	//dockerId := strings.TrimSpace(string(out))
	dockerId := "yubersvc-db-1"
	//fmt.Printf("Output: %s\n", dockerId)
	json2 := V(sh.Output("docker", "inspect", dockerId))
	//log.Printf("json: %s\n", json)
	//theId := gjson.Get(json, "0.NetworkSettings").String()
	//log.Printf("theId: %s\n", theId)
	host := "127.0.0.1"
	//x := gjson.Get(json, "0.NetworkSettings.Ports")
	//log.Printf("x: %+v\n", x)
	publishedPort := gjson.Get(json2, "0.NetworkSettings.Ports.5432/tcp.0.HostPort").Int()
	urlDb := url.URL{
		Scheme: "postgresql",
		Host:   fmt.Sprintf("%s:%d", host, publishedPort),
		User: url.UserPassword(
			gjson.Get(json, "services.ap.environment.DB_USER").String(),
			gjson.Get(json, "services.ap.environment.DB_PASSWORD").String(),
		),
		Path:    "/" + gjson.Get(json, "services.ap.environment.DB_DATABASE").String(),
		RawPath: "sslmode=disable",
	}
	V0(fmt.Fprintf(os.Stdout, "DB_URL=%s\n", urlDb.String()))
}

// Print (host string, port int, user, password, database string) prints text in the .env format that references the CDK configuration.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Env) Print(host string, port int, user, password, database string) {
	urlDb := url.URL{
		Scheme: "postgresql",
		Host:   fmt.Sprintf("%s:%d", host, port),
		User: url.UserPassword(
			user,
			password,
		),
		Path:     "/" + database,
		RawQuery: "sslmode=disable",
	}
	V0(fmt.Fprintf(os.Stdout, "DB_URL=%s\n", urlDb.String()))
}

var lintFns []any

// AddLintFn adds a function to the list of functions to lint something.
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func AddLintFn(fn any) {
	lintFns = append(lintFns, fn)
}

// Lint analyses.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Lint() error {
	mg.Deps(lintFns...)
	return nil
}

var genFns []any

// AddGenFn adds a function to the list of functions to generate something.
func AddGenFn(fn any) {
	genFns = append(genFns, fn)
}

// Gen generates binding codes.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Gen() error {
	// Docker build does not generate dependent files.
	if os.Getenv("NO_GEN") == "" && os.Getenv("NO_GENERATE") == "" {
		mg.Deps(genFns...)
	}
	return nil
}

// Generate is an alias for Gen.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Generate() error {
	return Gen()
}
