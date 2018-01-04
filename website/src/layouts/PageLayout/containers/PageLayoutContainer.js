import { connect } from 'react-redux'

import PageLayout from '../components/PageLayout'
import { logoutUser, loginUserSuccess } from '../../../store/authentication'

const mapDispatchToProps = {
    logout: () => logoutUser(),
    login: () => loginUserSuccess("token")
}
  
const mapStateToProps = (state) => ({
    token: state.authentication.token
})

export default connect(mapStateToProps, mapDispatchToProps)(PageLayout)