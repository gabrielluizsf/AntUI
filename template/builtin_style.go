package template

// builtinStyle paints exactly what the ready-made widgets paint — the same
// theme colours, the same built-in face, the same proportions — so the
// Builtin template is the look AntUI already had, at no cost. It exists to be
// the plainspoken neighbour of the futuristic one: same components, same
// events, same keyboard, different neon. Everything it draws is multiplied by
// [Scale], so the look survives a window that grows to any resolution.
type builtinStyle struct{}
