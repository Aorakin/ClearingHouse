```mermaid
---
title: POST /quota/organization - CreateOrganizationQuota
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: CreateOrganizationQuota(CreateOrganizationQuotaRequest)
    QuotaUsecase->>Repository: CreateOrgQuota()
    Repository-->>QuotaUsecase: organization quota
    QuotaUsecase-->>User: created organization quota
```

```mermaid
---
title: GET /quota/organization - GetOrganizationQuota
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: GetOrganizationQuota()
    QuotaUsecase->>Repository: GetOrganizationByRelationship()
    Repository-->>QuotaUsecase: organization quotas
    QuotaUsecase-->>User: organization quotas
```

```mermaid
---
title: GET /quota/organization/:org_id - GetOrganizationQuotasByOrgID
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: GetOrganizationQuotasByOrgID()
    QuotaUsecase->>Repository: GetOrganizationQuotasByOrgID()
    Repository-->>QuotaUsecase: organization quotas
    QuotaUsecase-->>User: organization quotas
```

```mermaid
---
title: PUT /quota/organization/:quota_id - UpdateOrganizationQuota
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: UpdateOrganizationQuota()
    QuotaUsecase->>Repository: GetOrgQuotaByID()
    QuotaUsecase->>Repository: UpdateOrganizationQuota()
    Repository-->>QuotaUsecase: updated organization quota
    QuotaUsecase-->>User: updated organization quota
```

```mermaid
---
title: DELETE /quota/organization/:quota_id - DeleteOrganizationQuota
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: DeleteOrganizationQuota()
    QuotaUsecase->>Repository: GetOrgQuotaByID()
    QuotaUsecase->>Repository: DeleteOrganizationQuota()
    Repository-->>QuotaUsecase: delete result
    QuotaUsecase-->>User: success
```

```mermaid
---
title: POST /quota/project - CreateProjectQuota
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: CreateProjectQuota(CreateProjectQuotaRequest)
    QuotaUsecase->>Repository: CreateProjectQuota()
    Repository-->>QuotaUsecase: project quota
    QuotaUsecase-->>User: created project quota
```

```mermaid
---
title: GET /quota/project/:project_id - GetProjectQuotas
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: GetProjectQuotas()
    QuotaUsecase->>Repository: GetProjectQuotaByProjectID()
    Repository-->>QuotaUsecase: project quotas
    QuotaUsecase-->>User: project quotas
```

```mermaid
---
title: GET /quota/project/:project_id/total - GetProjectQuotaTotal
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: GetProjectQuotaTotal()
    QuotaUsecase->>Repository: GetProjectQuotaTotalByType()
    Repository-->>QuotaUsecase: total quota
    QuotaUsecase-->>User: total quota
```

```mermaid
---
title: GET /quota/project/:project_id/namespaces - GetNamespaceQuotaInProject
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: GetNamespaceQuotaInProject()
    QuotaUsecase->>Repository: GetNamespaceQuotaInProject()
    Repository-->>QuotaUsecase: namespace quotas
    QuotaUsecase-->>User: namespace quotas
```

```mermaid
---
title: PUT /quota/project/:quota_id - UpdateProjectQuota
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: UpdateProjectQuota()
    QuotaUsecase->>Repository: GetProjectQuotaByID()
    QuotaUsecase->>Repository: UpdateProjectQuota()
    Repository-->>QuotaUsecase: updated project quota
    QuotaUsecase-->>User: updated project quota
```

```mermaid
---
title: POST /quota/project/internal - CreateInternalProjectQuota
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: CreateInternalProjectQuota(CreateInternalProjectQuotaRequest)
    QuotaUsecase->>Repository: CreateProjectQuota()
    Repository-->>QuotaUsecase: project quota
    QuotaUsecase-->>User: created project quota
```

```mermaid
---
title: PUT /quota/project/internal/:quota_id - UpdateInternalProjectQuota
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: UpdateInternalProjectQuota()
    QuotaUsecase->>Repository: GetProjectQuotaByID()
    QuotaUsecase->>Repository: UpdateInternalProjectQuota()
    Repository-->>QuotaUsecase: updated project quota
    QuotaUsecase-->>User: updated project quota
```

```mermaid
---
title: DELETE /quota/project/:quota_id - DeleteProjectQuota
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: DeleteProjectQuota()
    QuotaUsecase->>Repository: GetProjectQuotaByID()
    QuotaUsecase->>Repository: DeleteProjectQuota()
    Repository-->>QuotaUsecase: delete result
    QuotaUsecase-->>User: success
```

```mermaid
---
title: POST /quota/namespace - CreateNamespaceQuota
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: CreateNamespaceQuota(CreateNamespaceQuotaRequest)
    QuotaUsecase->>Repository: CreateNamespaceQuota()
    Repository-->>QuotaUsecase: namespace quota
    QuotaUsecase-->>User: created namespace quota
```

```mermaid
---
title: GET /quota/namespace/:namespace_id - GetNamespaceQuota
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: GetNamespaceQuota()
    QuotaUsecase->>Repository: GetNamespaceQuotaByNamespaceID()
    Repository-->>QuotaUsecase: namespace quotas
    QuotaUsecase-->>User: namespace quotas
```

```mermaid
---
title: PUT /quota/namespace/:quota_id - UpdateNamespaceQuota
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: UpdateNamespaceQuota()
    QuotaUsecase->>Repository: GetNamespaceQuotaByID()
    QuotaUsecase->>Repository: UpdateNamespaceQuota()
    Repository-->>QuotaUsecase: updated namespace quota
    QuotaUsecase-->>User: updated namespace quota
```

```mermaid
---
title: POST /quota/namespace/template - CreateNamespaceQuotaTemplate
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: CreateNamespaceQuotaTemplate(CreateNamespaceQuotaTemplateRequest)
    QuotaUsecase->>Repository: CreateNamespaceQuotaTemplate()
    Repository-->>QuotaUsecase: quota template
    QuotaUsecase-->>User: created quota template
```

```mermaid
---
title: GET /quota/namespace/template/:quota_template_id - GetNamespaceQuotaTemplate
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: GetNamespaceQuotaTemplate()
    QuotaUsecase->>Repository: GetNamespaceQuotaTemplateByID()
    Repository-->>QuotaUsecase: quota template
    QuotaUsecase-->>User: quota template
```

```mermaid
---
title: GET /quota/namespace/template/project/:project_id - GetNamespaceQuotaTemplatesByProjectID
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: GetNamespaceQuotaTemplatesByProjectID()
    QuotaUsecase->>Repository: GetNamespaceQuotaTemplatesByProjectID()
    Repository-->>QuotaUsecase: quota templates
    QuotaUsecase-->>User: quota templates
```

```mermaid
---
title: PUT /quota/namespace/template/:quota_template_id - UpdateNamespaceQuotaTemplate
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: UpdateNamespaceQuotaTemplate()
    QuotaUsecase->>Repository: GetNamespaceQuotaTemplateByID()
    QuotaUsecase->>Repository: UpdateNamespaceQuotaTemplate()
    Repository-->>QuotaUsecase: updated quota template
    QuotaUsecase-->>User: updated quota template
```

```mermaid
---
title: POST /quota/namespace/template/assign - AssignQuotaTemplateToNamespace
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: AssignQuotaTemplateToNamespace(AssignQuotaToNamespaceRequest)
    QuotaUsecase->>Repository: AssignQuotaTemplateToNamespace()
    Repository-->>QuotaUsecase: assign result
    QuotaUsecase-->>User: success
```

```mermaid
---
title: DELETE /quota/namespace/template/unassign/:namespace_id - UnassignQuotaTemplateFromNamespace
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: UnassignQuotaTemplateFromNamespace()
    QuotaUsecase->>Repository: UnassignQuotaTemplateFromNamespace()
    Repository-->>QuotaUsecase: unassign result
    QuotaUsecase-->>User: success
```

```mermaid
---
title: DELETE /quota/namespace/template/:quota_template_id - DeleteNamespaceQuotaTemplate
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: DeleteNamespaceQuotaTemplate()
    QuotaUsecase->>Repository: GetNamespaceQuotaTemplateByID()
    QuotaUsecase->>Repository: DeleteNamespaceQuotaTemplate()
    Repository-->>QuotaUsecase: delete result
    QuotaUsecase-->>User: success
```

```mermaid
---
title: DELETE /quota/namespace/:quota_id - DeleteNamespaceQuota
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: DeleteNamespaceQuota()
    QuotaUsecase->>Repository: GetNamespaceQuotaByID()
    QuotaUsecase->>Repository: DeleteNamespaceQuota()
    Repository-->>QuotaUsecase: delete result
    QuotaUsecase-->>User: success
```

```mermaid
---
title: GET /quota/:quota_id/usage/:namespace_id - GetUsage
---
sequenceDiagram
    actor User
    participant QuotaUsecase
    participant Repository

    User->>QuotaUsecase: GetUsage()
    QuotaUsecase->>Repository: IsAssigned()
    QuotaUsecase->>Repository: GetNamespaceUsageByType()
    Repository-->>QuotaUsecase: usage
    QuotaUsecase-->>User: usage
```
