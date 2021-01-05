import React, { useEffect, useState } from 'react';
import { Api } from '../api';
import { ClusterTable } from '../component/cluster-table';

export const ClustersPage: React.FunctionComponent = () => {
    let [state, setState] = useState([])
    useEffect(() => {
        Api.getClusters().then(
            res => setState(res)
        )
    }, [])
    return (
        <ClusterTable clusters={state}></ClusterTable>
    );
};