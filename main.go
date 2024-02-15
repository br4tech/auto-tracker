package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/base"
	"github.com/bluenviron/gortsplib/v4/pkg/format"
	"github.com/bluenviron/gortsplib/v4/pkg/format/rtph264"
	"github.com/bluenviron/mediacommon/pkg/codecs/h264"
	"github.com/otiai10/gosseract"
	"github.com/pion/rtp"
)

func main() {
	c := gortsplib.Client{}

	//URL da camera
	u, err := base.ParseURL("rtsp://admin:Guilherme@461@192.168.1.14:554/cam/realmonitor?channel=1&subtype=0")
	if err != nil {
		panic(err)
	}

	// Conectar-se ao servidor
	err = c.Start(u.Scheme, u.Host)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// Encontrar mídias disponíveis
	desc, _, err := c.Describe(u)
	if err != nil {
		panic(err)
	}

	// Encontre a mídia e o formato H264
	var forma *format.H264
	medi := desc.FindFormat(&forma)
	if medi == nil {
		panic("media not found")
	}

	// Configura RTP/H264 -> Decodificador H264
	rtpDec, err := forma.CreateDecoder()
	if err != nil {
		panic(err)
	}

	// Configurar H264 -> Decodificador de frames brutos
	frameDec := &h264Decoder{}
	err = frameDec.initialize()
	if err != nil {
		panic(err)
	}
	defer frameDec.close()

	// Se SPS e PPS estiverem presentes no SDP, envie-os para o decodificador
	if forma.SPS != nil {
		frameDec.decode(forma.SPS)
	}
	if forma.PPS != nil {
		frameDec.decode(forma.PPS)
	}

	// Configurar uma única mídia
	_, err = c.Setup(desc.BaseURL, medi, 0, 0)
	if err != nil {
		panic(err)
	}

	iframeReceived := false

	// Chamado quando um pacote RTP chega
	c.OnPacketRTP(medi, forma, func(pkt *rtp.Packet) {
		// Extrair unidades de acesso de pacotes RTP
		au, err := rtpDec.Decode(pkt)
		if err != nil {
			if err != rtph264.ErrNonStartingPacketAndNoPrevious && err != rtph264.ErrMorePacketsNeeded {
				log.Printf("ERR: %v", err)
			}
			return
		}

		//Esperar por um I-frame
		if !iframeReceived {
			if !h264.IDRPresent(au) {
				log.Printf("Esperar por um I-frame")
				return
			}
			iframeReceived = true
		}

		// Inicializar o cliente Tesseract
		clientOCR := gosseract.NewClient()
		// Definir o idioma para o reconhecimento de texto
		clientOCR.SetLanguage("por")

		for _, nalu := range au {
			// converter NALUs para RGBA frames
			img, err := frameDec.decode(nalu)
			if err != nil {
				panic(err)
			}

			// Espera por um frame para continuar o laco
			if img == nil {
				continue
			}

			// Definir nome do arquivo e criar arquivo temporario
			fname := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10) + ".jpg"
			f, err := os.CreateTemp("", fname)
			if err != nil {
				log.Fatalf("Erro ao criar o arquivo temporário: %v", err)
			}

			defer os.Remove(f.Name())
			defer f.Close()

			log.Println("saving", fname)

			// convert to jpeg
			err = jpeg.Encode(f, img, &jpeg.Options{
				Quality: 60,
			})

			if err != nil {
				log.Fatalf("Erro ao salvar a imagem: %v", err)
			}

			// Definir a imagem a ser processada pelo Tesseract-OCR
			err = clientOCR.SetImage(f.Name())
			if err != nil {
				log.Fatalf("Erro ao definir a imagem para o Tesseract-OCR: %v", err)
			}

			// Executar o OCR (reconhecimento óptico de caracteres)
			text, err := clientOCR.Text()
			if err != nil {
				log.Fatalf("Erro ao executar o OCR: %v", err)
			}

			// Exibir o texto reconhecido
			fmt.Println("Texto Reconhecido:", text)
		}
	})

	// start playing
	_, err = c.Play(nil)
	if err != nil {
		panic(err)
	}

	panic(c.Wait())
}
