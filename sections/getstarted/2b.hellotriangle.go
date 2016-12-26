package getstarted

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/raedatoui/learn-opengl-golang/utils"
	"github.com/raedatoui/learn-opengl-golang/sections"
)

type HelloSquare struct {
	sections.BaseSketch
	program       uint32
	vao, vbo, ebo uint32
}

func (hs *HelloSquare) Setup(w *glfw.Window, f *utils.Font) error {
	hs.Window = w
	hs.Font = f
	hs.Color = utils.RandColor()
	hs.Name = "2a. Hello Square"

	var vertexShader2 = `
	#version 330 core
	in vec3 vert;
	void main() {
		gl_Position = vec4(vert.x, vert.y, vert.z, 1.0);
	}` + "\x00"

	var fragShader2 = `
	#version 330 core
	out vec4 color;
	void main() {
		color = vec4(1.0f, 1.0f, 0.2f, 1.0f);
	}` + "\x00"

	var vertices = []float32{
		0.5, 0.5, 0.0, // Top Right
		0.5, -0.5, 0.0, // Bottom Right
		-0.5, -0.5, 0.0, // Bottom Left
		-0.5, 0.5, 0.0, // Top Left
	}

	var indices = []uint32{
		0, 1, 3, // First Triangle
		1, 2, 3, // Second Triangle
	}

	var err error
	hs.program, err = utils.BasicProgram(vertexShader2, fragShader2)
	if err != nil {
		return err
	}
	gl.UseProgram(hs.program)

	gl.GenVertexArrays(1, &hs.vao)
	gl.BindVertexArray(hs.vao)

	gl.GenBuffers(1, &hs.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, hs.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*utils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.GenBuffers(1, &hs.ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, hs.ebo)
	// seems like 4 works best here for the size of uint32
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(hs.program, gl.Str("vert\x00")))
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*utils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(vertAttrib)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return nil
}

func (hs *HelloSquare) Update() {

}

func (hs *HelloSquare) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(hs.Color.R, hs.Color.G, hs.Color.B, hs.Color.A)

	gl.UseProgram(hs.program)
	gl.BindVertexArray(hs.vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	gl.BindVertexArray(0)

	hs.Font.SetColor(0.0, 0.0, 0.0, 1.0)
	hs.Font.Printf(30, 30, 0.5, hs.Name)
}

func (hs *HelloSquare) Close() {
	gl.DeleteVertexArrays(1, &hs.vao)
	gl.DeleteBuffers(1, &hs.vbo)
	gl.DeleteBuffers(1, &hs.ebo)
	gl.UseProgram(0)
}

func (hs *HelloSquare) HandleKeyboard(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		hs.Window.SetShouldClose(true)
	}
}

func (hs *HelloSquare) HandleMousePosition(xpos, ypos float64) {

}

func (hs *HelloSquare) HandleScroll(xoff, yoff float64) {

}