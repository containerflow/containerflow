import React from 'react'
import { IndexLink, Link } from 'react-router'
import PropTypes from 'prop-types'
import './PageLayout.scss'

export const PageLayout = ({ children, token, logout, login }) => {
  return (
    <div className='container text-center'>
      <h1>ContainerFlow { token }</h1>
      <IndexLink to='/' activeClassName='page-layout__nav-item--active'>主页</IndexLink>
      {' · '}
      <Link to='/counter' activeClassName='page-layout__nav-item--active'>示例</Link>
      {' · '}
      { token ==='' || <a onClick={logout}>退出</a>}
      { token !=='' || <a onClick={login}>登录</a>}
      <div className='page-layout__viewport'>
        {children}
      </div>
    </div>
  )}

PageLayout.propTypes = {
  children: PropTypes.node,
  token: PropTypes.string,
  logout: PropTypes.func.isRequired,
  login: PropTypes.func.isRequired,
}

export default PageLayout
