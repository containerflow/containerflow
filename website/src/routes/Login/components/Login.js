import React from 'react'
import PropTypes from 'prop-types'

import { Card, Icon, Avatar,  Row, Col } from 'antd';

const { Meta } = Card;

export const Login = () => (
    <div>
        <Row>
            <Col span={12} order={1} >
            <Card
                style={{ width: 300 }}
                cover={<img alt="example" src="https://gw.alipayobjects.com/zos/rmsportal/JiqGstEfoWAOHiTxclqi.png" />}
            >
                <Meta
                title="Github"
                description="Login From Github"
                />
            </Card>
            </Col>
            <Col span={12} order={0} >
                <Card
                style={{ width: 300 }}
                cover={<img alt="example" src="https://gw.alipayobjects.com/zos/rmsportal/JiqGstEfoWAOHiTxclqi.png" />}
            >
                <Meta
                title="Gitlab"
                description="Login From Gitlab"
                />
            </Card>
            </Col>
        </Row>
    </div>
  )

export default Login