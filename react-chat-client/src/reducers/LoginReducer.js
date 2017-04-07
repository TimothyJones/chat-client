

export default function LoginReducer(state = {isLoggedIn: false},action) {
	console.log("LoginReducer called with ", state ," and action ", action)

	switch(action.type) {
		case "JOIN_ROOM":
			return {
				...state,
				isLoggedIn: true,
				name: action.name
			}
		default:
			return state;
	}
}