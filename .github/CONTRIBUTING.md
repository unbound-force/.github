# Contributing to Unbound Force

Thanks for joining the swarm. This guide applies to every Unbound Force repo
unless a repo has its own `CONTRIBUTING.md`. The full guide, including the
Speckit workflow in detail, is at
[unboundforce.dev/docs/contributing](https://unboundforce.dev/docs/contributing/).

By taking part, you agree to follow our
[Code of Conduct](CODE_OF_CONDUCT.md).

## Where to start

| Repo | What it is |
| ---- | ---------- |
| [unbound-force](https://github.com/unbound-force/unbound-force) | The hero personas and the `uf` CLI |
| [gaze](https://github.com/unbound-force/gaze) | Test quality analysis via side-effect detection |
| [dewey](https://github.com/unbound-force/dewey) | Knowledge retrieval for AI agents |
| [replicator](https://github.com/unbound-force/replicator) | Multi-agent coordination |
| [website](https://github.com/unbound-force/website) | unboundforce.dev |
| [.github](https://github.com/unbound-force/.github) | Org settings, teams and these files |

Each repo's `README.md` covers setup, and its `AGENTS.md`, where present,
describes the workflow for AI coding agents. Projects that are still
incubating live in [Zero Dot Force](https://github.com/zero-dot-force), the
junior team.

## Reporting bugs and asking for features

Open an issue on the repo it affects, and pick the issue form that fits.
Include steps to reproduce, what you expected and what happened. For a
feature, describe the problem you're trying to solve before the solution.

Security problems don't belong in public issues. See
[SECURITY.md](SECURITY.md).

## Making a change

1. For anything beyond a small fix, open an issue first so we can agree on
   the approach before you write the code.
2. Fork the repo and create a branch from `main`, named with a number
   prefix, for example `003-feature-name`.
3. Take non-trivial work through the
   [Speckit](https://github.com/github/spec-kit) workflow, from
   specification to implementation. The project's constitution
   (`.specify/memory/constitution.md`) has the final say.
4. Write commit messages in the
   [Conventional Commits](https://www.conventionalcommits.org/) format, for
   example `fix: correct broken link in team navigation`.
5. Open a pull request against `main` and fill in the template. The Divisor
   review council reviews it.

## Before a pull request can merge

The default branch of each code repo is protected:

- The repo's required CI checks must pass.
- At least one maintainer must approve. Where the repo has a `CODEOWNERS`
  file, that approval must come from a code owner.
- A new push after approval needs approving again.
- Force pushes to `main` are blocked.

## Licensing

Each repo states its license in its `LICENSE` file. By contributing, you
agree that your contributions are licensed under that repo's license.
