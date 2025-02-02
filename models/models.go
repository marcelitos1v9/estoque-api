package models

import "go.mongodb.org/mongo-driver/bson/primitive"
import "time"

type Produto struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Nome            string            `bson:"nome" json:"nome"`
    Descricao       string            `bson:"descricao" json:"descricao"`
    Preco           float64           `bson:"preco" json:"preco"`
    PrecoPromocional float64          `bson:"preco_promocional,omitempty" json:"preco_promocional,omitempty"`
    Estoque         int               `bson:"estoque" json:"estoque"`
    Categoria       string            `bson:"categoria" json:"categoria"`
    Fornecedor      string            `bson:"fornecedor" json:"fornecedor"`
    CodigoBarras    string            `bson:"codigo_barras" json:"codigo_barras"`
    DataCriacao     time.Time         `bson:"data_criacao" json:"data_criacao"`
    UltimaAtualizacao time.Time       `bson:"ultima_atualizacao" json:"ultima_atualizacao"`
    Status          string            `bson:"status" json:"status"` // ativo, inativo, em_promocao
    ImagemURL       string            `bson:"imagem_url,omitempty" json:"imagem_url,omitempty"`
    Tags            []string          `bson:"tags,omitempty" json:"tags,omitempty"`
}
