import React, { useEffect, useState } from 'react';
import { Api } from '../api';
import { ClusterTable } from '../component/cluster-table';
import { useParams } from 'react-router-dom'
import { Table } from 'react-bootstrap';
import { Service } from '../types/service';

export const ClusterDetailPage: React.FunctionComponent = () => {
    let { clusterIdEncoded } = useParams();
    let clusterId = decodeURIComponent(clusterIdEncoded)
    let [state, setState] = useState([])
    useEffect(() => {
        Api.getClusterDetails(clusterId).then(
            res => setState(res)
        )
    }, [])
    return (
        <>
            <h2>{clusterId}</h2>
            <Table>
                <thead>
                    <tr>
                        <th>Service Name</th>
                        <th>Launch Type</th>
                        <th>Desired Count</th>
                        <th>Running Count</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {state.map((service: Service) => 
                        <tr key={service.ServiceArn}>
                            <td>{service.Name}</td>
                            <td>{service.LaunchType}</td>
                            <td>{service.DesiredCount}</td>
                            <td>{service.RunningCount}</td>
                            <td>BUTTONS GO HERE</td>
                        </tr>
                    )}
                </tbody>
            </Table>
        </>
        
    );
};