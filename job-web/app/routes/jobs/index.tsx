import { Link, useLoaderData } from "remix";

import { getJobs } from "~/job";
import type { Job } from "~/job";

export const loader = async () => {
    return getJobs();
};
export default function Jobs() {
    const jobs = useLoaderData<Job[]>();
    return (
        <main>
            <h1>jobs</h1>
            <table>
                <tr>
                <th>job_key</th>
                <th>in_use</th>
                <th>created_on</th>
                <th>last_used_on</th>
                </tr>
                {jobs.map((job) => (
                    <tr>
                        <td>{job.job_key}</td>
                        <td>{job.in_use}</td>
                        <td>{job.created_on}</td>
                        <td>{job.last_used_on}</td>
                    </tr>
                ))}
            </table>
        </main>
    );
}
