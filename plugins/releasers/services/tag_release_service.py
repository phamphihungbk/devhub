from releasers import push_tag_to_repo
from releasers.models.response import ReleaseResponse


class TagReleaseService:
    def trigger(self, payload) -> ReleaseResponse:
        push_tag_to_repo(
            repo_url=payload.repo_url,
            tag=payload.tag,
            target=payload.target,
            message=payload.notes or f"Release {payload.name}",
        )

        return ReleaseResponse.from_payload(payload)
