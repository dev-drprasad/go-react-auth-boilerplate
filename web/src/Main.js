import React from "react";
import { Layout } from "antd";

const { Content } = Layout;

function Main({ children }) {
  return (
    <Layout>
      <Layout className="site-layout">
        <Content style={{ padding: "32px 16px 16px 16px" }} id="main">
          {children}
        </Content>
      </Layout>
    </Layout>
  );
}

export default Main;
