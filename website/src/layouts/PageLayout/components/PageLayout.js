import React from 'react'
import { Layout, Menu, Dropdown, Breadcrumb, Icon } from 'antd'
import { IndexLink, Link } from 'react-router'
import PropTypes from 'prop-types'
import './PageLayout.scss'

const { Header, Content, Footer } = Layout
const SubMenu = Menu.SubMenu
const MenuItemGroup = Menu.ItemGroup

export const PageLayout = ({ children, token, logout, login }) => {
  return (
    <Layout className="layout">
      <Header style={{ padding: '0' }}>
        <div className="logo"></div>
        <div className="cf-main-menu">
          <Menu
            theme="dark"
            mode="horizontal"
            defaultSelectedKeys={['1']}
            style={{ lineHeight: '64px' }}
          >
            <Menu.Item key="1"><IndexLink to='/' activeClassName='page-layout__nav-item--active'>Home</IndexLink></Menu.Item>
            <Menu.Item key="2"><Link to='/counter' activeClassName='page-layout__nav-item--active'>Samples</Link></Menu.Item>
          </Menu>
        </div>
        <div className="cf-namespaces">
            <Menu
              theme="dark"
              mode="horizontal"
              style={{ lineHeight: '64px' }}
            >
              <SubMenu key="sub4" title={<span><Icon type="setting" />
                {token !== '' || <span>Please Login</span>}
                {token === '' || <span>Welcome, cflow</span>}
              </span>}>
                {token !== '' ||<Menu.Item key="4" ><a onClick={login}>login</a></Menu.Item>}
                {token === '' || <Menu.Item key="5" ><a onClick={logout}>logout</a></Menu.Item>}
              </SubMenu>
            </Menu>
        </div>
      </Header>
      <Content style={{ padding: '0' }}>
        <Header style={{ background: '#fff', textAlign: 'center'}}>
        </Header>
        <div className='container text-center'>
          Token: { token }
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
