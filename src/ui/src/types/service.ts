export type Service = {
    ClusterArn: string,
    ServiceArn: string,
    Name: string,
    CreatedAt: string,
    LaunchType: string,
    Status: string,
    DesiredCount: number,
    RunningCount: number,
    PendingCount: number
}