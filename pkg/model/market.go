package model

type Market string

const DefaultMarket Market = "AU"

func (m *Market) String() string {
	return string(*m)
}
