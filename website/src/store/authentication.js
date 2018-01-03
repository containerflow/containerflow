// ------------------------------------
// Constants
// ------------------------------------
export const LOGIN_USER_SUCCESS = 'LOGIN_USER_SUCCESS'
export const LOGOUT_USER_SUCCESS = 'LOGOUT_USER_SUCCESS'

// ------------------------------------
// Actions
// ------------------------------------
export function loginUserSuccess(token){
    localStorage.setItem('token', token)
    return {
        type: LOGIN_USER_SUCCESS,
        payload: token
    }
}

export function logoutUser(){
    localStorage.clear()
    return {
        type: LOGOUT_USER_SUCCESS
    }
}

// ------------------------------------
// Action Handlers
// ------------------------------------
const ACTION_HANDLERS = {
    [LOGIN_USER_SUCCESS] : (state, action) => {
        return Object.assign({}, state, {
            token: action.payload
        })
    },
    [LOGOUT_USER_SUCCESS]: (state, action) =>{
        return Object.assign({}, state, {token: ""})
    }

}

// ------------------------------------
// Reducer
// ------------------------------------
const initialState = {token: ""}
export default function authenticationReducer(state=initialState, action) {
    const handler = ACTION_HANDLERS[action.type]
    return handler ? handler(state, action) : state
}