import React, { useEffect, useState } from 'react';
import {Cluster} from '../types/cluster'
import {Link} from 'react-router-dom'

export type IClusterRow = {
    cluster: Cluster
}


export const ClusterRow: React.FunctionComponent<IClusterRow> = (props) => {
    const [show, toggleShow] = useState(true);

    return (
        <>
            <tr>
                <td>
                    <Link to={"/cluster/" + encodeURIComponent(props.cluster.ClusterArn)}>{props.cluster.ClusterName}</Link>
                </td>
                <td>
                    {props.cluster.RunningTasksCount}
                </td>
                <td>
                    {props.cluster.ActiveServicesCount}
                </td>
                <td>
                    {props.cluster.RegisteredContainerInstancesCount}
                </td>
            </tr>
        </>
    );
};