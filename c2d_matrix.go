package steering_force

type Matrix struct {
	_11, _12, _13 float64
	_21, _22, _23 float64
	_31, _32, _33 float64
}

func NewMatrix() *Matrix {
	return &Matrix{
		0.0, 0.0, 0.0,
		0.0, 0.0, 0.0,
		0.0, 0.0, 0.0,
	}
}

type C2DMatrix struct {
	matrix *Matrix
}

func (c2dm *C2DMatrix) _11(val float64) {
	c2dm.matrix._11 = val
}

func (c2dm *C2DMatrix) _12(val float64) {
	c2dm.matrix._12 = val
}

func (c2dm *C2DMatrix) _13(val float64) {
	c2dm.matrix._13 = val
}

func (c2dm *C2DMatrix) _21(val float64) {
	c2dm.matrix._21 = val
}

func (c2dm *C2DMatrix) _22(val float64) {
	c2dm.matrix._22 = val
}

func (c2dm *C2DMatrix) _23(val float64) {
	c2dm.matrix._23 = val
}

func (c2dm *C2DMatrix) _31(val float64) {
	c2dm.matrix._31 = val
}

func (c2dm *C2DMatrix) _32(val float64) {
	c2dm.matrix._32 = val
}

func (c2dm *C2DMatrix) _33(val float64) {
	c2dm.matrix._33 = val
}

func (c2dm *C2DMatrix) MatrixMultiply(mIn *Matrix) {
	mat_temp := NewMatrix()

	//first row
	mat_temp._11 = (c2dm.matrix._11 * mIn._11) + (c2dm.matrix._12 * mIn._21) + (c2dm.matrix._13 * mIn._31)
	mat_temp._12 = (c2dm.matrix._11 * mIn._12) + (c2dm.matrix._12 * mIn._22) + (c2dm.matrix._13 * mIn._32)
	mat_temp._13 = (c2dm.matrix._11 * mIn._13) + (c2dm.matrix._12 * mIn._23) + (c2dm.matrix._13 * mIn._33)

	//second
	mat_temp._21 = (c2dm.matrix._21 * mIn._11) + (c2dm.matrix._22 * mIn._21) + (c2dm.matrix._23 * mIn._31)
	mat_temp._22 = (c2dm.matrix._21 * mIn._12) + (c2dm.matrix._22 * mIn._22) + (c2dm.matrix._23 * mIn._32)
	mat_temp._23 = (c2dm.matrix._21 * mIn._13) + (c2dm.matrix._22 * mIn._23) + (c2dm.matrix._23 * mIn._33)

	//third
	mat_temp._31 = (c2dm.matrix._31 * mIn._11) + (c2dm.matrix._32 * mIn._21) + (c2dm.matrix._33 * mIn._31)
	mat_temp._32 = (c2dm.matrix._31 * mIn._12) + (c2dm.matrix._32 * mIn._22) + (c2dm.matrix._33 * mIn._32)
	mat_temp._33 = (c2dm.matrix._31 * mIn._13) + (c2dm.matrix._32 * mIn._23) + (c2dm.matrix._33 * mIn._33)

	c2dm.matrix = mat_temp
}

func (c2dm *C2DMatrix) TransformVector2D(vPoint Vector2D) Vector2D {
	tempX := (c2dm.matrix._11 * vPoint.X) + (c2dm.matrix._21 * vPoint.Y) + (c2dm.matrix._31)
	tempY := (c2dm.matrix._12 * vPoint.X) + (c2dm.matrix._22 * vPoint.Y) + (c2dm.matrix._32)

	return Vector2D{tempX, tempY}
}

func (c2dm *C2DMatrix) Identity() {
	c2dm.matrix._11 = 1
	c2dm.matrix._12 = 0
	c2dm.matrix._13 = 0

	c2dm.matrix._21 = 0
	c2dm.matrix._22 = 1
	c2dm.matrix._23 = 0

	c2dm.matrix._31 = 0
	c2dm.matrix._32 = 0
	c2dm.matrix._33 = 1
}

func (c2dm *C2DMatrix) Rotate(fwd Vector2D, side Vector2D) {
	mat := NewMatrix()

	mat._11 = fwd.X
	mat._12 = fwd.Y
	mat._13 = 0

	mat._21 = side.X
	mat._22 = side.Y
	mat._23 = 0

	mat._31 = 0
	mat._32 = 0
	mat._33 = 1

	//and multiply
	c2dm.MatrixMultiply(mat)
}

func NewC2DMatrix() *C2DMatrix {
	c2dm := &C2DMatrix{matrix: NewMatrix()}
	c2dm.Identity()
	return c2dm
}
