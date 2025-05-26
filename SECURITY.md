# Security Policy for SolVM

The SolVM team and community take security vulnerabilities seriously. We appreciate your efforts to responsibly disclose your findings, and we will make every effort to acknowledge your contributions.

## Supported Versions

We are committed to providing security updates for the following versions of SolVM. Please ensure you are using a supported version to benefit from the latest security patches.

| Version      | Supported          | Notes                                                                 |
| :----------- | :----------------- | :-------------------------------------------------------------------- |
| `1.3.x`      | :white_check_mark: | **Current stable release line.** Actively receiving security updates. |

We encourage all users to run the latest stable version within the `1.1.x` series to ensure they have the most up-to-date features and security patches. Critical security fixes may be backported to the `1.0.x` series for a limited time after a new `1.1.x` release.

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

If you believe you have found a security vulnerability in SolVM, please report it to us privately by emailing:

**me@juliaklee.wtf**

Please include the following details in your report:

1.  **A clear description of the vulnerability:** Explain the type of vulnerability and its potential impact.
2.  **The version(s) of SolVM affected:** Be as specific as possible (e.g., `v1.1.0`, `v1.0.3`).
3.  **Steps to reproduce the vulnerability:** Provide a clear, step-by-step guide, including any necessary Lua scripts or configuration, that allows us to reproduce the issue.
4.  **Proof-of-Concept (PoC):** If possible, include a PoC that demonstrates the vulnerability.
5.  **Any potential mitigations or workarounds** you are aware of.
6.  **Your name or alias** for acknowledgment (if desired).

**What to Expect After Reporting:**

*   You will receive an **acknowledgment of your report within 48 hours**.
*   We will investigate the reported vulnerability and aim to provide an **initial assessment of its validity and severity within 7 business days**.
*   We will maintain an **open line of communication** with you throughout the investigation and remediation process, providing updates on our progress.
*   If the vulnerability is confirmed, we will work on a fix and plan for a coordinated disclosure.
*   We will credit you for your discovery (unless you prefer to remain anonymous) in the release notes or security advisories once the vulnerability is patched.

**Disclosure Policy:**

We aim to publicly disclose vulnerabilities once a fix is available and has been deployed, or in coordination with you if you have a specific disclosure timeline. We generally prefer to patch vulnerabilities before public disclosure to protect users.

We kindly ask that you **do not disclose the vulnerability publicly** until we have had a reasonable amount of time to address it.

Thank you for helping keep SolVM secure!
