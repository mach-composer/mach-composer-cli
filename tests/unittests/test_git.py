import pytest
from mach import git


@pytest.mark.parametrize("with_changes", (True, False))
def test_commit(mocker, with_changes):
    mock = mocker.patch(
        "subprocess.check_output", return_value="main.yml" if with_changes else None
    )
    git.commit("A commit message")

    assert mock.call_count == 2 if with_changes else 1
    call_args = mock.call_args_list
    assert call_args[0][0][0] == ["git", "status", "--short"]

    if with_changes:
        assert call_args[1][0][0] == ["git", "commit", "-m", "A commit message"]
