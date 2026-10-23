package dbus

// methodCall is the message that calls member on the object at path.
func methodCall(serial uint32, dest, path, iface, member, sig string, body []byte) []byte {
	return encodeMessage(typeMethodCall, serial, []headerField{
		{code: fieldPath, sig: 'o', str: path},
		{code: fieldInterface, sig: 's', str: iface},
		{code: fieldMember, sig: 's', str: member},
		{code: fieldDestination, sig: 's', str: dest},
	}, sig, body)
}
