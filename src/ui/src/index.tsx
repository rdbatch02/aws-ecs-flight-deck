
import React, {  } from 'react';
import ReactDOM from 'react-dom';

// Importing the Bootstrap CSS
import 'bootstrap/dist/css/bootstrap.min.css';

import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import {
    BrowserRouter as Router,
    Switch,
    Route  } from "react-router-dom";
import { ClustersPage } from './page/clusters';
import { ClusterDetailPage } from './page/cluster-detail';

const App = () => {
    return (
        <Router>
            <Container className="p-3">
            <Row>
                <Col><h1>AWS ECS Flight Deck</h1></Col>
            </Row>
            
            <Switch>
                <Route path="/cluster/:clusterIdEncoded">
                    <ClusterDetailPage />
                </Route>
                <Route path="/">
                    <ClustersPage />
                </Route>
            </Switch>
            
        </Container>
        </Router>
        
    )

};

const mountNode = document.getElementById('root');
ReactDOM.render(<App />, mountNode);