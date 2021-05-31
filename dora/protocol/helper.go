package protocol

func (m *Pkg) Clone() *Pkg {
	head := *m.Head
	mc := &Pkg{
		Head:    &head,
		Payload: nil,
	}

	if len(mc.Head.Meta) == 0 {
		mc.Head.Meta = make(map[string]string)
	}

	return mc
}
