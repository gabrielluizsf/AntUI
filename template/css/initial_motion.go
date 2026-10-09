package css

func initialMotion(st *Style, prop string) {
	switch prop {
	case "transition":
		st.TransitionProps = nil
		st.TransitionDurs = nil
		st.TransitionTims = nil
		st.TransitionDels = nil
		st.TransitionNone = false
	case "transition-property":
		st.TransitionProps = nil
		st.TransitionNone = false
	case "transition-duration":
		st.TransitionDurs = nil
	case "transition-timing-function":
		st.TransitionTims = nil
	case "transition-delay":
		st.TransitionDels = nil
	case "animation":
		st.AnimationNames = nil
		st.AnimationDurs = nil
		st.AnimationTims = nil
		st.AnimationDels = nil
		st.AnimationIters = nil
		st.AnimationDirs = nil
		st.AnimationFills = nil
	case "animation-name":
		st.AnimationNames = nil
	case "animation-duration":
		st.AnimationDurs = nil
	case "animation-timing-function":
		st.AnimationTims = nil
	case "animation-delay":
		st.AnimationDels = nil
	case "animation-iteration-count":
		st.AnimationIters = nil
	case "animation-direction":
		st.AnimationDirs = nil
	case "animation-fill-mode":
		st.AnimationFills = nil
	}
}
