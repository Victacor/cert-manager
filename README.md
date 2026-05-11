# cert-manager

`cert-manager` es un controlador de certificados minimalista escrito en Go que ayuda a administrar un almacén PKI local.

## ¿Qué hace?

- Lee certificados y claves desde una estructura de carpetas `pki/`.
- Empareja cada certificado con su clave privada usando el `modulus` RSA.
- Muestra en pantalla una tabla con:
  - nombre del certificado
  - nombre de la clave
  - fecha de expiración
  - fecha de firma
- Colorea la fecha de expiración cuando el certificado tiene 60 días o menos para caducar.

## Estructura esperada

El proyecto asume una estructura de PKI similar a esta:

```text
pki/
  certs/
    *.crt
  keys/
    *.key
```

## Configuración y ruta PKI

- En esta versión, la ruta por defecto se encuentra en `./examples/pki`.
- El diseño del proyecto apunta a soportar una variable de entorno como `PKI_PATH` para usar una ruta personalizada.
- Si `PKI_PATH` no está definida, el binario buscaría por defecto en la ruta local del ejecutable.

## Uso

### Ejecutar el proyecto

```bash
go run ./cmd/main.go list
```

### Construir el binario

```bash
go build -o cert-manager ./cmd
./cert-manager list
```

## Salida esperada

El comando `list` imprime una tabla con los certificados válidos y sus claves emparejadas. La columna de expiración se muestra en rojo cuando faltan 60 días o menos para el vencimiento.

## Cómo funciona internamente

1. Busca archivos `*.crt` en `pki/certs`.
2. Busca archivos `*.key` en `pki/keys`.
3. Decodifica cada certificado y clave PEM.
4. Parsea los certificados X.509 y las claves RSA.
5. Extrae la fecha de firma (`NotBefore`) y la fecha de expiración (`NotAfter`).
6. Compara los módulos RSA para emparejar certificado y clave.
7. Renderiza una tabla usando `github.com/olekukonko/tablewriter`.

## Dependencias

- `github.com/spf13/cobra` para el CLI.
- `github.com/olekukonko/tablewriter` para visualizar tablas en la terminal.

## Mejoras sugeridas

El proyecto está preparado para ampliarse. Entre las funcionalidades que pueden añadirse:

- Comando para inspeccionar un certificado específico.
- Comando para renovar un CSR y generar una nueva firma cuando la clave privada está disponible.
- Mejor manejo de errores y soporte para más formatos de clave.

## Estado actual

Esta versión realiza únicamente el listado y emparejado de certificados y claves. El CLI está pensado para crecer en funciones de inspección y renovación.
