import React from 'react'
import { Layout, Menu, Breadcrumb } from 'antd'
import { IndexLink, Link } from 'react-router'
import PropTypes from 'prop-types'
import './PageLayout.scss'

const { Header, Content, Footer } = Layout;

export const PageLayout = ({ children, token, logout, login }) => {
  return (
    <Layout className="layout">
      <Header>
        <div className="logo"></div>
        <Menu
          theme="dark"
          mode="horizontal"
          defaultSelectedKeys={['2']}
          style={{ lineHeight: '64px' }}
        >
          <Menu.Item key="1"><IndexLink to='/' activeClassName='page-layout__nav-item--active'>主页</IndexLink></Menu.Item>
          <Menu.Item key="2"><Link to='/counter' activeClassName='page-layout__nav-item--active'>示例</Link></Menu.Item>
          <Menu.Item key="3">nav 3</Menu.Item>
        </Menu>
      </Header>
      <Content style={{ padding: '0 50px' }}>
        <div className='container text-center'>
          <h1>ContainerFlow { token }</h1>
          { token ==='' || <a onClick={logout}>退出</a>}
          { token !=='' || <a onClick={login}>登录</a>}
          <div className='page-layout__viewport'>
            {children}
          </div>
        </div>
      </Content>
      <Footer style={{ textAlign: 'center' }}>
        Ant Design ©2016 Created by Ant UED
      </Footer>
    </Layout>
  )}

PageLayout.propTypes = {
  children: PropTypes.node,
  token: PropTypes.string,
  logout: PropTypes.func.isRequired,
  login: PropTypes.func.isRequired,
}

export default PageLayout
