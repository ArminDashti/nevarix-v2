import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export function MonitoringDashboard() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Monitoring</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-muted-foreground">Fleet metrics will appear here after agents connect.</p>
      </CardContent>
    </Card>
  );
}
