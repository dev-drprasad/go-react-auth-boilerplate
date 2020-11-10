import React from "react";
import { Layout, Menu } from "antd";

const { Content, Sider } = Layout;

function Main({ children }) {
  return (
    <Layout className="site-layout" style={{ height: "100%" }}>
      <Sider width={200}>
        <Menu mode="inline" style={{ height: "100%" }}>
          <Menu.Item key="1">option1</Menu.Item>
        </Menu>
      </Sider>
      <Layout>
        <Content style={{ padding: "32px 16px 16px 16px" }} id="main">
          {children}
        </Content>
      </Layout>
    </Layout>
  );
}

export default Main;
