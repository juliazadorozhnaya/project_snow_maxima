package main

import (
	"fmt"
	"math"
	"sync"

	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	//graphics
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Scan", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure the vertex and fragment shaders
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 1, 60000.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	vColor := gl.GetUniformLocation(program, gl.Str("vColor\x00"))

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(triangles)*4, gl.Ptr(triangles), gl.DYNAMIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 3*4, 0)

	var vao2 uint32
	gl.GenVertexArrays(1, &vao2)
	gl.BindVertexArray(vao2)
	var vbo_plane uint32
	gl.GenBuffers(1, &vbo_plane)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo_plane)
	gl.BufferData(gl.ARRAY_BUFFER, len(dots)*4, gl.Ptr(dots), gl.DYNAMIC_DRAW)
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, 3*4, 0)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.PointSize(5)
	//gl.LineWidth(5)

	angle := 0.0
	previousTime := glfw.GetTime()

	var cubeMX sync.Mutex

	pcd := NewPointCloud(GetScanData(5304712))
	fmt.Println("Count of dots:", len(pcd.d))

	update_scene := func(edges []Edge) {
		ct := []float32{}
		for _, e := range edges {
			ct = append(ct, float32(e.b.x), float32(e.b.y), float32(e.b.z),
				float32(e.e.x), float32(e.e.y), float32(e.e.z))
		}
		cubeMX.Lock()
		triangles = ct
		cubeMX.Unlock()
	}
	//points array
	dots := []float32{}
	for _, d := range pcd.d {
		dots = append(dots, float32(d.x), float32(d.y), float32(d.z))
	}
	//code for triangulation
	go pcd.Triangulate(800000, update_scene)
	//triangulatedData := pcd.Triangulate(80000000, update_scene)
	//fmt.Println("Volume =", triangulatedData.GetVolume())

	//camera config to see all points
	minX := pcd.d[0].x
	maxX := pcd.d[len(pcd.d)-1].x
	if minX < 0 {
		minX = -minX
	}
	if maxX < 0 {
		maxX = -maxX
	}
	var pcdWidth int
	if minX > maxX {
		pcdWidth = minX
	} else {
		pcdWidth = maxX
	}
	camera := mgl32.LookAtV(mgl32.Vec3{float32(math.Tan(math.Pi*75/180)) * float32(pcdWidth) / 2, 0, 4 * float32(pcd.d[0].z)}, mgl32.Vec3{0, 0, float32(pcd.d[0].z)}, mgl32.Vec3{0, 0, 1})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	//go update_scene()
	for !window.ShouldClose() {
		cubeMX.Lock()
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(triangles)*4, gl.Ptr(triangles), gl.DYNAMIC_DRAW)
		cubeMX.Unlock()

		gl.BindBuffer(gl.ARRAY_BUFFER, vbo_plane)
		gl.BufferData(gl.ARRAY_BUFFER, len(dots)*4, gl.Ptr(dots), gl.DYNAMIC_DRAW)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		angle = (angle + elapsed) * 1
		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 0, 1})

		// Render
		gl.UseProgram(program)

		gl.Uniform3f(vColor, 0.0, 0.0, 1.0)
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		gl.BindVertexArray(vao)
		//gl.TRIANGLE_STRIP
		gl.DrawArrays(gl.LINES, 0, int32(len(triangles)/3))

		gl.Uniform3f(vColor, 0, 1.0, 1.0)
		gl.BindVertexArray(vao2)
		gl.DrawArrays(gl.POINTS, 0, int32(len(dots)/3))

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

var vertexShader = `
#version 410
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec2 vertTexCoord;
out vec2 fragTexCoord;
out vec4 vfcolor;
void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
}
` + "\x00"

var fragmentShader = `
#version 410
uniform sampler2D tex;
uniform vec3 vColor;
in vec4 vfcolor;
in vec2 fragTexCoord;
out vec4 outputColor;
void main() {
    outputColor = vec4(vColor, 0.5);
}
` + "\x00"

var triangles = []float32{
	0, 0.5, 0, // top
	-0.5, -0.5, 0, // left
	0.5, -0.5, 0, // right
}

var dots = []float32{
	-50, 288, -25,
}
