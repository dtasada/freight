package src

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	opt "github.com/moznion/go-optional"
)

type Object struct {
	Model rl.Model
}

func NewObjectFromFile(modelPath string, texturePath opt.Option[string], shader opt.Option[rl.Shader]) Object {
	model := rl.LoadModel(modelPath)

	if texturePath.IsSome() {
		model.Materials.Maps.Texture = rl.LoadTexture(texturePath.Unwrap())
	}

	if shader.IsSome() {
		model.Materials.Shader = shader.Unwrap()
	}

	return Object{
		model,
	}
}

func NewObjectFromMesh(mesh rl.Mesh, texturePath opt.Option[string], shader opt.Option[rl.Shader]) Object {
	model := rl.LoadModelFromMesh(mesh)

	if texturePath.IsSome() {
		model.Materials.Maps.Texture = rl.LoadTexture(texturePath.Unwrap())
	}

	if shader.IsSome() {
		model.Materials.Shader = shader.Unwrap()
	}

	return Object{
		model,
	}
}
