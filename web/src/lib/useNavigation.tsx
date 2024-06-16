import {
	createContext,
	useContext,
	createEffect,
	onCleanup,
	createSignal, Setter,
} from 'solid-js'
import { useBackButton } from './useBackButton'
import { useLocation, useNavigate } from '@solidjs/router'

interface NavigationContext {
	navigateBack: () => void,
	setPrevLocation: Setter<string>
}

const Navigation = createContext<NavigationContext>({} as NavigationContext)

export function NavigationProvider(props: { children: any }) {
	const backButton = useBackButton()

	const location = useLocation()

	const [prevLocation, setPrevLocation] = createSignal<string>(location.pathname)

	const navigate = useNavigate()

	const navigateBack = () => {
		if (location.pathname === prevLocation()) {
			navigate('/')
		} else if (prevLocation() !== '') {
			navigate(-1)
		} else {
			setPrevLocation(location.pathname)
			navigate('/')
		}
	}

	createEffect(() => {
		backButton.hide()
		if (location.pathname !== '/') {
			backButton.setVisible()
			backButton.onClick(navigateBack)
		}
	})

	onCleanup(() => {
		backButton.hide()
		backButton.offClick(navigateBack)
	})

	const value: NavigationContext = {
		navigateBack,
		setPrevLocation,
	}

	return (
		<Navigation.Provider value={value}>{props.children}</Navigation.Provider>
	)
}

export function useNavigation() {
	return useContext(Navigation)
}
