import React, { useEffect, useState } from 'react';
import { Table } from 'react-bootstrap';
import Col from 'react-bootstrap/esm/Col';
import Row from 'react-bootstrap/esm/Row';
import { Cluster } from '../types/cluster'
import { ClusterRow } from './cluster-row';

export type IClusterTable = {
    clusters: Cluster[]
}

export const ClusterTable: React.FunctionComponent<IClusterTable> = (props) => {
    return (
        <>
            <Table striped bordered hover>
                <thead>
                    <tr>
                        <th>Cluster Name</th>
                        <th>Running Task Count</th>
                        <th>Active Service Count</th>
                        <th>Container Instance Count</th>
                    </tr>
                    
                </thead>


                <tbody>
                    {props.clusters.map(cluster =>
                        <ClusterRow cluster={cluster as Cluster} key={cluster.ClusterArn}></ClusterRow>
                    )}
                </tbody>
            </Table>
        </>
    );
};