package gitea

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"code.gitea.io/sdk/gitea"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	repoOwner                    string = "username"
	repoName                     string = "name"
	repoDescription              string = "description"
	repoPrivateFlag              string = "private"
	repoIssueLabels              string = "issue_labels"
	repoAutoInit                 string = "auto_init"
	repoTemplate                 string = "repo_template"
	repoGitignores               string = "gitignores"
	repoLicense                  string = "license"
	repoReadme                   string = "readme"
	repoDefaultBranch            string = "default_branch"
	repoWebsite                  string = "website"
	repoIssues                   string = "has_issues"
	repoWiki                     string = "has_wiki"
	repoPrs                      string = "has_pull_requests"
	repoProjects                 string = "has_projects"
	repoIgnoreWhitespace         string = "ignore_whitespace_conflicts"
	repoAllowMerge               string = "allow_merge_commits"
	repoAllowRebase              string = "allow_rebase"
	repoAllowRebaseMerge         string = "allow_rebase_explicit"
	repoAllowSquash              string = "allow_squash_merge"
	repoArchived                 string = "archived"
	repoAllowManualMerge         string = "allow_manual_merge"
	repoAutodetectManualMerge    string = "autodetect_manual_merge"
	repoMirror                   string = "mirror"
	repoArchiveOnDestroy         string = "archive_on_destroy"
	migrationCloneAddresse       string = "migration_clone_addresse"
	migrationCloneAddress        string = "migration_clone_address"
	migrationService             string = "migration_service"
	migrationServiceAuthName     string = "migration_service_auth_username"
	migrationServiceAuthPassword string = "migration_service_auth_password"
	migrationServiceAuthToken    string = "migration_service_auth_token"
	migrationMilestones          string = "migration_milestones"
	migrationReleases            string = "migration_releases"
	migrationIssueLabels         string = "migration_issue_labels"
	migrationMirrorInterval      string = "migration_mirror_interval"
	migrationLFS                 string = "migration_lfs"
	migrationLFSEndpoint         string = "migration_lfs_endpoint"
)

func searchUserByName(c *gitea.Client, name string) (res *gitea.User, err error) {
	page := 1

	for {
		users, _, err := c.AdminListUsers(gitea.AdminListUsersOptions{
			ListOptions: gitea.ListOptions{
				Page:     page,
				PageSize: 50,
			},
		})
		if err != nil {
			return nil, err
		}

		if len(users) == 0 {
			return nil, fmt.Errorf("User with name %s could not be found", name)
		}

		for _, user := range users {
			if user.UserName == name {
				return user, nil
			}
		}

		page += 1
	}
}

func resourceRepoRead(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*gitea.Client)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	var resp *gitea.Response

	if err != nil {
		return err
	}

	repo, resp, err := client.GetRepoByID(id)

	if err != nil {
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		} else {
			return err
		}
	}

	err = setRepoResourceData(repo, d)

	return
}

func resourceRepoCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*gitea.Client)

	var repo *gitea.Repository
	var resp *gitea.Response
	var orgRepo, hasAdmin bool

	_, resp, err = client.GetOrg(d.Get(repoOwner).(string))

	if resp.StatusCode == 404 {
		_, err := searchUserByName(client, d.Get(repoOwner).(string))
		if err != nil {
			if strings.Contains(err.Error(), "could not be found") {
				return errors.New(fmt.Sprintf("Creation of repository cound not proceed as owner %s is not present in gitea", d.Get(repoOwner).(string)))
			}
			tflog.Warn(context.Background(), "Error query for users. Assuming missing permissions and proceding with user permissions")
			hasAdmin = false
		} else {
			hasAdmin = true
		}
		orgRepo = false
	} else {
		orgRepo = true
	}

	var cloneAddr string
	if d.Get(migrationCloneAddresse).(string) != "" {
		cloneAddr = d.Get(migrationCloneAddresse).(string)
	} else {
		cloneAddr = d.Get(migrationCloneAddress).(string)
	}

	if cloneAddr != "" {

		opts := gitea.MigrateRepoOption{
			RepoName:       d.Get(repoName).(string),
			RepoOwner:      d.Get(repoOwner).(string),
			CloneAddr:      cloneAddr,
			Service:        gitea.GitServiceType(d.Get(migrationService).(string)),
			Mirror:         d.Get(repoMirror).(bool),
			Private:        d.Get(repoPrivateFlag).(bool),
			Description:    d.Get(repoDescription).(string),
			Wiki:           d.Get(repoWiki).(bool),
			Milestones:     d.Get(migrationMilestones).(bool),
			Labels:         d.Get(migrationIssueLabels).(bool),
			Issues:         d.Get(repoIssues).(bool),
			PullRequests:   d.Get(repoPrs).(bool),
			Releases:       d.Get(migrationReleases).(bool),
			MirrorInterval: d.Get(migrationMirrorInterval).(string),
			LFS:            d.Get(migrationLFS).(bool),
			LFSEndpoint:    d.Get(migrationLFSEndpoint).(string),
		}

		if d.Get(migrationServiceAuthName).(string) != "" {
			opts.AuthUsername = d.Get(migrationServiceAuthName).(string)
		}
		if d.Get(migrationServiceAuthPassword).(string) != "" {
			opts.AuthPassword = d.Get(migrationServiceAuthPassword).(string)
		}
		if d.Get(migrationServiceAuthToken).(string) != "" {
			opts.AuthToken = d.Get(migrationServiceAuthToken).(string)
		}

		repo, _, err = client.MigrateRepo(opts)

	} else {
		opts := gitea.CreateRepoOption{
			Name:          d.Get(repoName).(string),
			Description:   d.Get(repoDescription).(string),
			Private:       d.Get(repoPrivateFlag).(bool),
			IssueLabels:   d.Get(repoIssueLabels).(string),
			AutoInit:      d.Get(repoAutoInit).(bool),
			Template:      d.Get(repoTemplate).(bool),
			Gitignores:    d.Get(repoGitignores).(string),
			License:       d.Get(repoLicense).(string),
			Readme:        d.Get(repoReadme).(string),
			DefaultBranch: d.Get(repoDefaultBranch).(string),
			TrustModel:    "default",
		}

		if orgRepo {
			repo, _, err = client.CreateOrgRepo(d.Get(repoOwner).(string), opts)
		} else {
			if hasAdmin {
				repo, _, err = client.AdminCreateRepo(d.Get(repoOwner).(string), opts)
			} else {
				repo, _, err = client.CreateRepo(opts)
			}
		}
	}

	if err != nil {
		return err
	}

	err = setRepoResourceData(repo, d)

	return
}

func resourceRepoUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*gitea.Client)

	var repo *gitea.Repository

	var name string = d.Get(repoName).(string)
	var description string = d.Get(repoDescription).(string)
	var website string = d.Get(repoWebsite).(string)
	var private bool = d.Get(repoPrivateFlag).(bool)
	var template bool = d.Get(repoTemplate).(bool)
	var hasIssues bool = d.Get(repoIssues).(bool)
	var hasWiki bool = d.Get(repoWiki).(bool)
	var defaultBranch string = d.Get(repoDefaultBranch).(string)
	var hasPRs bool = d.Get(repoPrs).(bool)
	var hasProjects bool = d.Get(repoProjects).(bool)
	var ignoreWhitespaceConflicts bool = d.Get(repoIgnoreWhitespace).(bool)
	var allowMerge bool = d.Get(repoAllowMerge).(bool)
	var allowRebase bool = d.Get(repoAllowRebase).(bool)
	var allowRebaseMerge bool = d.Get(repoAllowRebaseMerge).(bool)
	var allowSquash bool = d.Get(repoAllowSquash).(bool)
	var allowManualMerge bool = d.Get(repoAllowManualMerge).(bool)
	var autodetectManualMerge bool = d.Get(repoAutodetectManualMerge).(bool)

	opts := gitea.EditRepoOption{
		Name:                      &name,
		Description:               &description,
		Website:                   &website,
		Private:                   &private,
		Template:                  &template,
		HasIssues:                 &hasIssues,
		HasWiki:                   &hasWiki,
		DefaultBranch:             &defaultBranch,
		HasPullRequests:           &hasPRs,
		HasProjects:               &hasProjects,
		IgnoreWhitespaceConflicts: &ignoreWhitespaceConflicts,
		AllowMerge:                &allowMerge,
		AllowRebase:               &allowRebase,
		AllowRebaseMerge:          &allowRebaseMerge,
		AllowSquash:               &allowSquash,
		AllowManualMerge:          &allowManualMerge,
		AutodetectManualMerge:     &autodetectManualMerge,
	}

	if d.Get(repoMirror).(bool) {
		var mirrorInterval string = d.Get(migrationMirrorInterval).(string)
		opts.MirrorInterval = &mirrorInterval
	} else {
		var archived bool = d.Get(repoArchived).(bool)
		opts.Archived = &archived
	}

	repo, _, err = client.EditRepo(d.Get(repoOwner).(string), d.Get(repoName).(string), opts)

	if err != nil {
		return err
	}
	err = setRepoResourceData(repo, d)

	return

}

func resourceRepoDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*gitea.Client)

	archiveOnDestroy := d.Get(repoArchiveOnDestroy).(bool)
	archived := d.Get(repoArchived).(bool)
	owner := d.Get(repoOwner).(string)
	name := d.Get(repoName).(string)

	if archiveOnDestroy {
		if archived {
			log.Printf("[DEBUG] Repository already archived, nothing to do on delete: %s/%s", owner, name)
			err = nil
			return err
		} else {
			log.Printf("[DEBUG] Archiving repository on delete: %s/%s", owner, name)
			err = archiveRepo(d, client)
			return err
		}
	} else {
		log.Printf("[DEBUG] Deleting repository: %s/%s", owner, repoName)
		err = deleteRepo(d, client)
		return err
	}
}

func archiveRepo(d *schema.ResourceData, client *gitea.Client) (err error) {
	if err := d.Set("archived", true); err != nil {
		return err
	}

	var name string = d.Get(repoName).(string)
	var description string = d.Get(repoDescription).(string)
	var website string = d.Get(repoWebsite).(string)
	var private bool = d.Get(repoPrivateFlag).(bool)
	var template bool = d.Get(repoTemplate).(bool)
	var hasIssues bool = d.Get(repoIssues).(bool)
	var hasWiki bool = d.Get(repoWiki).(bool)
	var defaultBranch string = d.Get(repoDefaultBranch).(string)
	var hasPRs bool = d.Get(repoPrs).(bool)
	var hasProjects bool = d.Get(repoProjects).(bool)
	var ignoreWhitespaceConflicts bool = d.Get(repoIgnoreWhitespace).(bool)
	var allowMerge bool = d.Get(repoAllowMerge).(bool)
	var allowRebase bool = d.Get(repoAllowRebase).(bool)
	var allowRebaseMerge bool = d.Get(repoAllowRebaseMerge).(bool)
	var allowSquash bool = d.Get(repoAllowSquash).(bool)
	var allowManualMerge bool = d.Get(repoAllowManualMerge).(bool)
	var autodetectManualMerge bool = d.Get(repoAutodetectManualMerge).(bool)
	var archived bool = d.Get(repoArchived).(bool)

	opts := gitea.EditRepoOption{
		Name:                      &name,
		Description:               &description,
		Website:                   &website,
		Private:                   &private,
		Template:                  &template,
		HasIssues:                 &hasIssues,
		HasWiki:                   &hasWiki,
		DefaultBranch:             &defaultBranch,
		HasPullRequests:           &hasPRs,
		HasProjects:               &hasProjects,
		IgnoreWhitespaceConflicts: &ignoreWhitespaceConflicts,
		AllowMerge:                &allowMerge,
		AllowRebase:               &allowRebase,
		AllowRebaseMerge:          &allowRebaseMerge,
		AllowSquash:               &allowSquash,
		AllowManualMerge:          &allowManualMerge,
		AutodetectManualMerge:     &autodetectManualMerge,
		Archived:                  &archived,
	}

	_, _, err = client.EditRepo(d.Get(repoOwner).(string), d.Get(repoName).(string), opts)
	return err
}

func deleteRepo(d *schema.ResourceData, client *gitea.Client) (err error) {
	_, err = client.DeleteRepo(d.Get(repoOwner).(string), d.Get(repoName).(string))
	return err
}

func setRepoResourceData(repo *gitea.Repository, d *schema.ResourceData) (err error) {
	d.SetId(fmt.Sprintf("%d", repo.ID))
	d.Set("username", repo.Owner.UserName)
	d.Set("name", repo.Name)
	d.Set("description", repo.Description)
	d.Set("full_name", repo.FullName)
	d.Set("private", repo.Private)
	d.Set("fork", repo.Fork)
	d.Set("mirror", repo.Mirror)
	d.Set("size", repo.Size)
	d.Set("html_url", repo.HTMLURL)
	d.Set("ssh_url", repo.SSHURL)
	d.Set("clone_url", repo.CloneURL)
	d.Set("website", repo.Website)
	d.Set("stars", repo.Stars)
	d.Set("forks", repo.Forks)
	d.Set("watchers", repo.Watchers)
	d.Set("open_issue_count", repo.OpenIssues)
	d.Set("default_branch", repo.DefaultBranch)
	d.Set("created", repo.Created.String())
	d.Set("updated", repo.Updated.String())
	d.Set("permission_admin", repo.Permissions.Admin)
	d.Set("permission_push", repo.Permissions.Push)
	d.Set("permission_pull", repo.Permissions.Pull)

	return
}

func resourceGiteaRepository() *schema.Resource {
	return &schema.Resource{
		Read:   resourceRepoRead,
		Create: resourceRepoCreate,
		Update: resourceRepoUpdate,
		Delete: resourceRepoDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Owner of the repository",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Name of the repository",
			},
			"auto_init": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Default:     true,
				Description: "Flag if the repository should be initiated with the configured values",
			},
			"repo_template": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
			"issue_labels": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "Default",
				Description: "The Issue Label configuration to be used in this repository.\n" +
					"Need to exist in the gitea instance",
			},
			"gitignores": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "",
				Description: "A specific gitignore that should be commited to the repository" +
					"on creation if `auto_init` is set to `true`\n" +
					"Need to exist in the gitea instance",
			},
			"license": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "",
				Description: "The license under which the source code of this repository should be.\n" +
					"Need to exist in the gitea instance",
			},
			"readme": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Default:     "",
				Description: "The description of the repository.",
			},
			"private": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Default:     true,
				Description: "Flag if the repository should be private or not.",
			},
			"default_branch": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Default:     "main",
				Description: "The default branch of the repository. Defaults to `main`",
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permission_admin": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"permission_push": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"permission_pull": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"website": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Default:     "",
				Description: "A link to a website with more information.",
			},
			"has_issues": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Default:     true,
				Description: "A flag if the repository should have issue management enabled or not.",
			},
			"has_wiki": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Default:     true,
				Description: "A flag if the repository should have the native wiki enabled or not.",
			},
			"has_pull_requests": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Default:     true,
				Description: "A flag if the repository should acceppt pull requests or not.",
			},
			"has_projects": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Default:     true,
				Description: "A flag if the repository should have the native project management enabled or not.",
			},
			"ignore_whitespace_conflicts": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  true,
			},
			"allow_merge_commits": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  true,
			},
			"allow_rebase": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  true,
			},
			"allow_rebase_explicit": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  true,
			},
			"allow_squash_merge": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  true,
			},
			"archived": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
			"archive_on_destroy": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Default:     false,
				Description: "Set to 'true' to archive the repository instead of deleting on destroy.",
			},
			"allow_manual_merge": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  true,
			},
			"autodetect_manual_merge": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  true,
			},
			"mirror": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  false,
			},
			"migration_clone_addresse": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "DEPRECATED in favor of `migration_clone_address`",
			},
			"migration_clone_address": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				ForceNew: true,
			},
			"migration_service": {
				Type:        schema.TypeString,
				Required:    false,
				ForceNew:    true,
				Optional:    true,
				Description: "git/github/gitlab/gitea/gogs",
			},
			"migration_service_auth_username": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "",
			},
			"migration_service_auth_password": {
				Type:      schema.TypeString,
				Required:  false,
				Optional:  true,
				Sensitive: true,
				Default:   "",
			},
			"migration_service_auth_token": {
				Type:      schema.TypeString,
				Required:  false,
				Optional:  true,
				Sensitive: true,
				Default:   "",
			},
			"migration_milestones": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  true,
			},
			"migration_releases": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  true,
			},
			"migration_issue_labels": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
				Default:  true,
			},
			"migration_mirror_interval": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Default:     "8h0m0s",
				Description: "valid time units are 'h', 'm', 's'. 0 to disable automatic sync",
			},
			"migration_lfs": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
			},
			"migration_lfs_endpoint": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "",
			},
			"clone_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"html_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssh_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Description: "`gitea_repository` manages a gitea repository.\n\n" +
			"Per default this repository will be initializiled with the provided configuration (gitignore, License etc.).\n" +
			"If the `username` property is set to a organisation name, the provider will try to look if this organisation exists " +
			"and create the repository under the organisation scope.\n\n" +
			"Repository migrations have some properties that are not available to regular repositories. These are all prefixed with `migration_`.",
	}
}
