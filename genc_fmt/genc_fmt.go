package genc_fmt

import (
	"fmt"
	"log"
	"regexp"
	"strings"

)

//  structs and enums : -------------------------------------------------------------- (section)  //
type Directive struct {

	d_type string
	name   string
	args   string
}

type Column struct {
	name string
	values string
}

type Table struct {
	name string
	cols []Column
}
//  (section) -------------------------------------------------------------- : structs and enums  //

func extractDirective(line string) Directive {
	// fmt.Println(line)
	directive_finding_re, err := regexp.Compile(`(@|(\[.*\])|(\(.*\)))`)
	if err != nil {
		log.Fatalf("Failure with generating regex: %v", err)
	}
	directive := directive_finding_re.ReplaceAllString(line, "")
	name := ""
	args := ""
	is_name_present, err := regexp.MatchString(`\[.*\]`, line)
	if err != nil {
		log.Fatalf("Error identifying if name is present for the directive => %v", err)
	}
	if is_name_present {
		name_finding_re, err := regexp.Compile(`(@.*\[)|(\](\([^)]*\))?)`)
		if err != nil {
			log.Fatalf("Failure with generating regex: %v", err)
		}
		name = name_finding_re.ReplaceAllString(line, "")
	}
	is_args_present, err := regexp.MatchString(`\([^)]*\)`, line)
	if err != nil {
		log.Fatalf("Error identifying if args are present for the directive => %v", err)
	}
	if is_args_present {
		args_finding_re, err := regexp.Compile(`(@.*(\[.*\])?\()|(\))`)
		if err != nil {
			log.Fatalf("Failure with generating regex: %v", err)
		}
		args = args_finding_re.ReplaceAllString(line, "")
	}

	return Directive{
		directive,
		name,
		args,
	}
}


func parseATable(lines []string, ln *int) Table {
	for i := range() {
	}
}

func ParseGenc(content []byte) {
	input := string(content[:])
	lines := strings.Split(input, "\n")

	for ln := range len(lines){
		is_directive_start, err := regexp.MatchString(`^@.*(\[.*\])??(\(.*\))??`, lines[ln])
		if err != nil {
			log.Fatalf("Failure with finding directives: %v", err)
		}
		if is_directive_start {
			directive := extractDirective(lines[ln])
			fmt.Printf("directives =>  %v\n\n", directive)

			switch directive.d_type {
			case "table":
				table := parseATable(lines, &ln)
				fmt.Println(table)

			case "enum":

			case "enum2string":

			case "gen_func_types":

			case "gen_func_globals":

			case "gen_custom_template":

			default:
				log.Fatalf("not found a viable directive type => %s", directive.d_type)
			}
		}
	}
}
