package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type Certificate struct {
	name            string
	date_expiration string
	date_signature  string
	mod             string
}

type Key struct {
	name string
	mod  string
}

func main() {

	var cmdPrint = &cobra.Command{
		Use:   "list [string to print]",
		Short: "Ayuda para cobra",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			files, keys, err := obtainFiles()
			if err != nil {
				fmt.Println(err)
			}
			row_certs := getinfofiles(files, keys)
			formatTable(row_certs)
			//color_date("2026-05-8")
		},
	}
	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdPrint)
	rootCmd.Execute()
}

func color_date(date_s string) string {
	red := "\033[31m"
	green := "\033[32m"
	reset := "\033[0m"

	layout := "2006-01-02"

	expirationDate, _ := time.Parse(layout, date_s)

	now := time.Now()

	diff := expirationDate.Sub(now)
	limit := int(diff / 24)

	if limit <= 60 {
		//fmt.Println(red + date_s + reset)
		colored_string := (red + date_s + reset)
		return colored_string
	} else {
		//fmt.Println(green + date_s + reset)
		colored_string := (green + date_s + reset)
		return colored_string
	}
}

func obtainFiles() ([]string, []string, error) {
	pki_path := os.Getenv("PKI_DIR")
	if pki_path == "" {
		pki_path = "./pki"
	}

	certificates_exp := (pki_path + "/certs/*.crt")
	keys_exp := (pki_path + "/keys/*.key")

	files, err := filepath.Glob(certificates_exp)
	if err != nil {
		return nil, nil, err
	}
	keys, err := filepath.Glob(keys_exp)
	if err != nil {
		return nil, nil, err
	}
	return files, keys, nil
}

func getinfofiles(files_path []string, keys_path []string) [][]string {
	var certificates []Certificate
	var keys []Key
	var row_certs [][]string

	for _, v := range files_path {
		cert_name := filepath.Base(v)
		cert_data, cert_d_err := os.ReadFile(v)
		if cert_d_err != nil {
			fmt.Println("Valor nil al leer el cert")
			os.Exit(1)
		}
		cert_decoded, _ := pem.Decode(cert_data)
		if cert_decoded == nil {
			fmt.Println("Error decodificando el cert")
			os.Exit(1)
		}
		cert_human_readable, cert_p_err := x509.ParseCertificate(cert_decoded.Bytes)
		if cert_p_err != nil {
			fmt.Println("Error parseando el certificado")
			os.Exit(1)
		}
		cert_date_expiration := cert_human_readable.NotAfter.Format("2006-01-02")
		if cert_data == nil {
			fmt.Println("No se puede obtener la fecha del crt")
			os.Exit(1)
		}
		cert_date_signature := cert_human_readable.NotBefore.Format("2006-01-02")
		cert_pub := cert_human_readable.PublicKey.(*rsa.PublicKey)
		cert_n := cert_pub.N
		if cert_n == nil {
			fmt.Println("No se puede obtener el modulus del crt")
			os.Exit(1)
		}
		cert_n_str := fmt.Sprintf("%x", cert_n)

		certificates = append(certificates, Certificate{name: cert_name, date_expiration: cert_date_expiration, date_signature: cert_date_signature, mod: cert_n_str})
	}

	for _, v := range keys_path {
		key_name := filepath.Base(v)
		key_data, key_d_err := os.ReadFile(v)
		if key_d_err != nil {
			fmt.Println("Valor empty al leer el key")
			os.Exit(1)
		}
		key_decoded, _ := pem.Decode(key_data)
		if key_decoded == nil {
			fmt.Println("Error decodificando la key")
			os.Exit(1)
		}
		key_human_readble, key_p_err := x509.ParsePKCS8PrivateKey(key_decoded.Bytes)
		if key_p_err != nil {
			key_pkcs1, _ := x509.ParsePKCS1PrivateKey(key_decoded.Bytes)
			key_n_pkcs1 := key_pkcs1.N
			key_n_str_pkcs1 := fmt.Sprintf("%x", key_n_pkcs1)

			keys = append(keys, Key{name: key_name, mod: key_n_str_pkcs1})
			continue
		}
		key_rsa := key_human_readble.(*rsa.PrivateKey)
		key_n := key_rsa.N
		key_n_str := fmt.Sprintf("%x", key_n)

		keys = append(keys, Key{name: key_name, mod: key_n_str})
	}

	for _, v1 := range certificates {
		for _, v2 := range keys {
			if v1.mod == v2.mod {
				row_certs = append(row_certs, []string{v1.name, v2.name, color_date(v1.date_expiration), v1.date_signature}) //v1.date_expiration, v1.date_signature
			}
		}
	}

	return row_certs
}

func formatTable(row_certs [][]string) {

	data := [][]string{
		{"Certificado", "Key", "Valided", "Firma"},
	}

	data = append(data, row_certs...)
	table := tablewriter.NewWriter(os.Stdout)
	table.Header(data[0])
	table.Bulk(data[1:])
	table.Render()
}
