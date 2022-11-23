package core

import (
	"fmt"
	"strings"
)

func ProtocolMap(m *WorldMap) string {
	sb := new(strings.Builder)
	sb.WriteString("START MAP\n")

	for yRow := 0; yRow < m.YHeight(); yRow++ {
		for xCol := 0; xCol < m.XWidth(); xCol++ {
			sb.WriteByte(m.Cell(xCol, yRow).Type())
		}
		sb.WriteByte('\n')
	}

	sb.WriteString("END MAP\n")
	return sb.String()
}

func ProtocolPlayer(m *WorldMap) string {
	sb := new(strings.Builder)
	sb.WriteString("START PLAYER\n")

	for _, player := range m.Players() {
		sb.WriteString("PlayerID:")
		sb.WriteString(fmt.Sprintf("%d", player.PlayerID()))

		sb.WriteString("|Name:")
		sb.WriteString(player.Name())

		sb.WriteString("|Color:")
		sb.WriteString(player.Color())

		sb.WriteString("|Position:")
		sb.WriteString(fmt.Sprintf("%.6f,%.6f", player.Position().X(), player.Position().Y()))

		sb.WriteString("|Velocity:")
		sb.WriteString(fmt.Sprintf("%.6f,%.6f", player.Velocity().X(), player.Velocity().Y()))

		sb.WriteString("|Acceleration:")
		sb.WriteString(fmt.Sprintf("%.6f,%.6f", player.Acceleration().X(), player.Acceleration().Y()))

		sb.WriteString("|Score:")
		sb.WriteString(fmt.Sprintf("%d", player.Score()))

		sb.WriteString("|Angle:")
		sb.WriteString(fmt.Sprintf("%.6f", player.Angle()))

		sb.WriteString("|TouchingCells:")
		for _, tc := range player.TouchingCells() {
			sb.WriteString(fmt.Sprintf("%d,%d;", tc.XCol(), tc.YRow()))
		}

		sb.WriteString("|IsAlive:")
		sb.WriteString(fmt.Sprintf("%v", player.IsAlive()))

		sb.WriteByte('\n')
	}

	sb.WriteString("END PLAYER\n")
	return sb.String()
}

func ProtocolStatus(m *WorldMap) string {
	sb := new(strings.Builder)
	sb.WriteString("START STATUS\n")

	iteration, endtime, maxUpdateTime := m.Stats()
	sb.WriteString(fmt.Sprintf("Iteration:%d\n", iteration))
	sb.WriteString(fmt.Sprintf("Endtime:%d\n", endtime))
	sb.WriteString(fmt.Sprintf("MaxUpdateTime:%v\n", maxUpdateTime))
	sb.WriteString(fmt.Sprintf("MaxPlayers:%d\n", m.MaxPlayers()))

	sb.WriteString("END STATUS\n")
	return sb.String()
}
