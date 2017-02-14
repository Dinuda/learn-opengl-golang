package getstarted

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/raedatoui/glutils"
	"github.com/raedatoui/learn-opengl-golang/sections"
)

type HelloCamera struct {
	sections.BaseSketch
	shader               uint32
	vao, vbo, ebo        uint32
	texture1, texture2   uint32
	transform            mgl32.Mat4
	cubePositions        []mgl32.Mat4
	camera               glutils.Camera
	deltaTime, lastFrame float64
	w, a, s, d           bool
	lastX, lastY         float64
	firstMouse           bool
	rotationAxis         mgl32.Vec3
}

func (hc *HelloCamera) InitGL() error {
	hc.Name = "7. Camera (use WSDA and mouse)"

	var err error
	hc.shader, err = glutils.Shader("_assets/getting_started/6.coordinates/coordinate.vs",
		"_assets/getting_started/6.coordinates/coordinate.frag", "")
	if err != nil {
		return err
	}

	gl.UseProgram(hc.shader)

	vertices := []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
	}

	hc.rotationAxis = mgl32.Vec3{1.0, 0.3, 0.5}.Normalize()

	hc.cubePositions = []mgl32.Mat4{
		mgl32.Translate3D(0.0, 0.0, 0.0),
		mgl32.Translate3D(2.0, 5.0, -15.0),
		mgl32.Translate3D(-1.5, -2.2, -2.5),
		mgl32.Translate3D(-3.8, -2.0, -12.3),
		mgl32.Translate3D(2.4, -0.4, -3.5),
		mgl32.Translate3D(-1.7, 3.0, -7.5),
		mgl32.Translate3D(1.3, -2.0, -2.5),
		mgl32.Translate3D(1.5, 2.0, -2.5),
		mgl32.Translate3D(1.5, 0.2, -1.5),
		mgl32.Translate3D(-1.3, 1.0, -1.5),
	}

	// ====================
	// camera
	// ====================
	hc.camera = glutils.NewCamera(
		mgl32.Vec3{0.0, 0.0, 3.0},
		mgl32.Vec3{0.0, 1.0, 3.0},
		glutils.YAW, glutils.PITCH,
	)
	hc.lastX = 400
	hc.lastY = 300
	hc.firstMouse = true

	gl.GenVertexArrays(1, &hc.vao)
	gl.GenBuffers(1, &hc.vbo)
	gl.GenBuffers(1, &hc.ebo)

	gl.BindVertexArray(hc.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, hc.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*glutils.GL_FLOAT32_SIZE, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*glutils.GL_FLOAT32_SIZE, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// TexCoord attribute
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 5*glutils.GL_FLOAT32_SIZE, gl.PtrOffset(3*glutils.GL_FLOAT32_SIZE))
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0) // Unbind VAO

	// Texture 1
	if tex, err := glutils.NewTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "_assets/images/container.png"); err != nil {
		return err
	} else {
		hc.texture1 = tex
	}

	// Texture 2
	if tex, err := glutils.NewTexture(gl.REPEAT, gl.REPEAT, gl.LINEAR, gl.LINEAR, "_assets/images/awesomeface.png"); err != nil {
		return err
	} else {
		hc.texture2 = tex
	}

	return nil
}

func (hc *HelloCamera) Update() {
	// Set frame time
	currentFrame := glfw.GetTime()
	hc.deltaTime = currentFrame - hc.lastFrame
	hc.lastFrame = currentFrame
	if hc.w {
		hc.camera.ProcessKeyboard(glutils.FORWARD, hc.deltaTime)
	}
	if hc.s {
		hc.camera.ProcessKeyboard(glutils.BACKWARD, hc.deltaTime)
	}
	if hc.a {
		hc.camera.ProcessKeyboard(glutils.LEFT, hc.deltaTime)
	}
	if hc.d {
		hc.camera.ProcessKeyboard(glutils.RIGHT, hc.deltaTime)
	}
}

func (hc *HelloCamera) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(hc.Color32.R, hc.Color32.G, hc.Color32.B, hc.Color32.A)

	// Bind Textures using texture units
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, hc.texture1)
	loc1 := gl.GetUniformLocation(hc.shader, gl.Str("ourTexture1\x00"))
	gl.Uniform1i(loc1, 0)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, hc.texture2)
	loc2 := gl.GetUniformLocation(hc.shader, gl.Str("ourTexture2\x00"))
	gl.Uniform1i(loc2, 1)

	// Activate shader
	gl.UseProgram(hc.shader)

	// Create camera transformations
	view := hc.camera.GetViewMatrix()
	projection := mgl32.Perspective(float32(hc.camera.Zoom), sections.RATIO, 0.1, 1000.0)

	// Get their uniform location
	modelLoc := gl.GetUniformLocation(hc.shader, gl.Str("model\x00"))
	viewLoc := gl.GetUniformLocation(hc.shader, gl.Str("view\x00"))
	projLoc := gl.GetUniformLocation(hc.shader, gl.Str("projection\x00"))
	// Pass the matrices to the shader
	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
	// Note: currently we set the projection matrix each frame,
	// but since the projection matrix rarely changes it's often best practice to set it outside the main loop only once.
	gl.UniformMatrix4fv(projLoc, 1, false, &projection[0])

	// Draw container
	gl.BindVertexArray(hc.vao)

	for i := 0; i < 10; i++ {
		// Calculate the model matrix for each object and pass it to shader before drawing
		model := hc.cubePositions[i]
		angle := float32(glfw.GetTime()) * float32(i+1)
		model = model.Mul4(mgl32.HomogRotate3D(angle, hc.rotationAxis))
		gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
		gl.DrawArrays(gl.TRIANGLES, 0, 36)
	}
	gl.BindVertexArray(0)
}

func (hc *HelloCamera) Close() {
	gl.DeleteVertexArrays(1, &hc.vao)
	gl.DeleteBuffers(1, &hc.vbo)
	gl.DeleteBuffers(1, &hc.ebo)
	gl.DeleteProgram(hc.shader)
}

func (hc *HelloCamera) HandleKeyboard(k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey, keys map[glfw.Key]bool) {
	hc.w = keys[glfw.KeyW]
	hc.a = keys[glfw.KeyA]
	hc.s = keys[glfw.KeyS]
	hc.d = keys[glfw.KeyD]
}

func (hc *HelloCamera) HandleMousePosition(xpos, ypos float64) {
	if hc.firstMouse {
		hc.lastX = xpos
		hc.lastY = ypos
		hc.firstMouse = false
	}

	xoffset := xpos - hc.lastX
	yoffset := hc.lastY - ypos // Reversed since y-coordinates go from bottom to left

	hc.lastX = xpos
	hc.lastY = ypos

	hc.camera.ProcessMouseMovement(xoffset, yoffset, true)
}

func (hc *HelloCamera) HandleScroll(xoff, yoff float64) {
	hc.camera.ProcessMouseScroll(yoff)
}
