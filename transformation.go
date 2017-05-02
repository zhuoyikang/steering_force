package steering_force

import "fmt"

// 将世界坐标转换为局部坐标
func PointToLocalSpace(point Vector2D, AgentHeading Vector2D,
	AgentSide Vector2D, AgentPosition Vector2D) Vector2D {
	matTransform := NewC2DMatrix()

	Tx := -AgentPosition.Dot(AgentHeading)
	Ty := -AgentPosition.Dot(AgentSide)

	matTransform._11(AgentHeading.X)
	matTransform._12(AgentSide.X)
	matTransform._21(AgentHeading.Y)
	matTransform._22(AgentSide.Y)
	matTransform._31(Tx)
	matTransform._32(Ty)

	return matTransform.TransformVector2D(point)
}

func VectorToWorldSpace(point Vector2D, AgentHeading Vector2D,
	AgentSide Vector2D, AgentPosition Vector2D) Vector2D {

	matTransform := NewC2DMatrix()
	//rotate
	matTransform.Rotate(AgentHeading, AgentSide)

	fmt.Printf("VectorToWorldSpace %v\n", *matTransform.matrix)

	//now transform the vertices
	return matTransform.TransformVector2D(point)
}
