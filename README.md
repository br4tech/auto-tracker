# Auto Traker de Monitoramento de Placas de Veículos

## Descrição
Este é um projeto de monitoramento de placas de veículos desenvolvido em Go (Golang), usando PostgreSQL como banco de dados e GORM como ORM. Ele usa a biblioteca Go Tesseract (gosseract) para reconhecimento óptico de caracteres (OCR) em imagens de placas de veículos.

## Pré-requisitos
- Go (Golang)
- PostgreSQL
- GORM
- Tesseract OCR
- Gosseract

## Instalação
1. Clone este repositório.
2. Instale as dependências do Go:

```
go mod tidy
```

3. Instale e configure o PostgreSQL.
4. Instale o Tesseract OCR e configure-o de acordo com as instruções do [Tesseract OCR GitHub](https://github.com/tesseract-ocr/tesseract).
E necessario instalar essa lib para o tesseract funcionar
[TESSDOC](https://github.com/tesseract-ocr/tessdoc?tab=readme-ov-file)
5. Instale a biblioteca Go Tesseract (gosseract):

```
go get github.com/otiai10/gosseract/v2
sudo apt-get install -y libavformat-dev libswscale-dev gcc pkg-config
sudo add-apt-repository ppa:alex-p/tesseract-ocr-devel
sudo apt-get install libtesseract-dev
sudo apt-get install tesseract-ocr-eng
```

## Configuração
1. Configure as variáveis de ambiente necessárias, como credenciais de banco de dados e configurações do Tesseract OCR.

## Uso
1. Execute o aplicativo:

``` 
go run main.go
```


## Contribuição
Contribuições são bem-vindas! Sinta-se à vontade para enviar pull requests ou abrir issues.

## Licença
[MIT License](LICENSE)



