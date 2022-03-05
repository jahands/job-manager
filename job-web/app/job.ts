export type Job = {
  job_key: string;
  in_use: boolean;
  created_on: string;
  last_used_on: string;
};

type JobsResult = {
  result: Job[];
}

export async function getJobs() {
  const res = await fetch("https://job-manager-staging-aa4e.up.railway.app/v1/test/jobs/?api_key=x");
  if (!res.ok) {
    throw new Error(res.statusText)
  }
  const posts: Job[] = (await res.json() as JobsResult).result;
  return posts;
}
