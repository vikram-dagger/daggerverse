from typing import Annotated, Self
import dagger
from dagger import Container, dag, Directory, DefaultPath, Doc, File, Secret, function, object_type, ReturnType
import requests
import re

@object_type
class GithubApi:

    @function
    async def create_comment(
        self,
        repository: Annotated[str, Doc("The owner and repository name")],
        ref: Annotated[str, Doc("The ref name")],
        body: Annotated[str, Doc("The comment body")],
        token: Annotated[Secret, Doc("GitHub API token")],
    ) -> str:
        """Adds a comment to the PR"""
        pr_number = int(re.search(r"(\d+)", ref).group(1))
        plaintext = await token.plaintext()

        url = f"https://api.github.com/repos/{repository}/issues/{pr_number}/comments"
        headers = {
            "Authorization": f"Bearer {plaintext}",
            "Accept": "application/vnd.github+json"
        }
        data = {
            "body": body
        }
        response = requests.post(url, headers=headers, json=data)

        if response.status_code == 201:
            return f"{response.json()['html_url']}"
        else:
            raise Exception(f"Failed to post comment: {response.status_code} - {response.text}")

    @function
    async def create_pr(
        self,
        repository: Annotated[str, Doc("The owner and repository name")],
        ref: Annotated[str, Doc("The ref name")],
        diff_file: Annotated[File, Doc("The diff file")],
        token: Annotated[Secret, Doc("GitHub API token")],
    ) -> str:
        """Creates a new PR with the changes"""
        plaintext = await token.plaintext()
        pr_number = int(re.search(r"(\d+)", ref).group(1))
        new_branch = f"patch-from-pr-{pr_number}"
        remote_url = f"https://${{GITHUB_TOKEN}}@github.com/{repository}.git"
        diff = await diff_file.contents()

        await (
            dag
            .container()
            .from_("alpine/git")
            .with_new_file("/tmp/a.diff", f"{diff}")
            .with_workdir("/app")
            .with_env_variable("GITHUB_TOKEN", plaintext)
            .with_exec(["git", "init"])
            .with_exec(["git", "config", "user.name", "Dagger Agent"])
            .with_exec(["git", "config", "user.email", "vikram@dagger.io"])
            .with_exec(["sh", "-c", "git remote add origin " + remote_url])
            .with_exec(["git", "fetch", "origin", f"pull/{pr_number}/head:{new_branch}"])
            .with_exec(["git", "checkout", new_branch])
            .with_exec(["git", "apply", "/tmp/a.diff"])
            .with_exec(["git", "add", "."])
            .with_exec(["git", "commit", "-m", f"Fixes PR #{pr_number}"])
            .with_exec(["git", "push", "--set-upstream", "origin", new_branch])
            .stdout()
        )

        headers = {
            "Authorization": f"Bearer {plaintext}",
            "Accept": "application/vnd.github+json"
        }
        pr_url = f"https://api.github.com/repos/{repository}/pulls/{pr_number}"
        pr_response = requests.get(pr_url, headers=headers)

        if pr_response.status_code != 200:
            raise Exception(f"Failed to fetch original PR: {pr_response.text}")

        pr_data = pr_response.json()
        base_branch = pr_data["head"]["ref"]

        create_pr_url = f"https://api.github.com/repos/{repository}/pulls"
        head_user = repository.split("/")[0]
        head = f"{head_user}:{new_branch}"

        payload = {
            "title": f"Automated follow-up to PR #{pr_number}",
            "body": f"This PR fixes PR #{pr_number} using `{new_branch}`.",
            "head": head,
            "base": base_branch
        }

        create_response = requests.post(create_pr_url, headers=headers, json=payload)
        if create_response.status_code != 201:
            raise Exception(f"Failed to create new PR: {create_response.text}")

        new_pr = create_response.json()
        return f"{new_pr['html_url']}"
