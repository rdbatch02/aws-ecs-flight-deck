
import React, { useEffect, useState } from 'react';
import ReactDOM from 'react-dom';

// Importing the Bootstrap CSS
import 'bootstrap/dist/css/bootstrap.min.css';

import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Toast from 'react-bootstrap/Toast';
import Container from 'react-bootstrap/Container';
import Button from 'react-bootstrap/Button';
import { Api } from './api';

const ExampleToast: React.FunctionComponent = ({ children }) => {
    const [show, toggleShow] = useState(true);

    return (
        <>
            {!show && <Button onClick={() => toggleShow(true)}>Show Toast</Button>}
            <Toast show={show} onClose={() => toggleShow(false)}>
                <Toast.Header>
                    <strong className="mr-auto">React-Bootstrap</strong>
                </Toast.Header>
                <Toast.Body>{children}</Toast.Body>
            </Toast>
        </>
    );
};

const App = () => {
    let [state, setState] = useState([])
    useEffect(() => {
        Api.getClusters().then(
            res => setState(res)
        )
    }, [])

    return (
        <Container className="p-3">
            <Row>
                <Col><h1>AWS ECS Flight Deck</h1></Col>
            </Row>
            <Row>
                {state.map( cluster => <div>{JSON.stringify(cluster)}</div>)}
            </Row>
        </Container>
    )

};

const mountNode = document.getElementById('root');
ReactDOM.render(<App />, mountNode);