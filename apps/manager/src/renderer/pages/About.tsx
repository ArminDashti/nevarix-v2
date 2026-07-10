import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export function About() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>About</CardTitle>
      </CardHeader>
      <CardContent>
        <p>Nevarix Manager v0.1.0</p>
        <p className="text-sm text-muted-foreground">Server monitoring and management platform</p>
      </CardContent>
    </Card>
  );
}
