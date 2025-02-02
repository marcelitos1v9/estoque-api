package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nome           string            `bson:"nome" json:"nome"`
	Email          string            `bson:"email" json:"email"`
	Senha          string            `bson:"senha" json:"-"` // O "-" impede que a senha seja enviada no JSON
	Role           string            `bson:"role" json:"role"` // admin, manager, user
	DataCriacao    time.Time         `bson:"data_criacao" json:"data_criacao"`
	UltimoAcesso   time.Time         `bson:"ultimo_acesso" json:"ultimo_acesso"`
	Ativo          bool              `bson:"ativo" json:"ativo"`
} 