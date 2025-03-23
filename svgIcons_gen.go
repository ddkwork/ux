package ux

import (
	"embed"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/golibrary/stream"
	"github.com/ddkwork/ux/giosvg"
)

func Svg2Icon(b []byte) *giosvg.Icon {
	return giosvg.NewIcon(mylog.Check2(giosvg.NewVector(b)))
}

func svgCallback(value []byte) []byte {
	if strings.Contains(string(value), "fill=\"none\"") {
		// return []byte(strings.Replace(string(value), "fill=\"none\"", "fill=\"white\"", 1))
	}
	return []byte(strings.Replace(string(value), "<path", "<path fill=\"white\"", 1))
}

// 取色
// https://products.eptimize.app/zh/color-convert/rgb-to-rgba
//
//go:embed resources/images/*.svg
var images embed.FS

var (
	svgEmbedFileMap                = stream.ReadEmbedFileMap(images, "resources/images")
	SvgIconBookmark                = Svg2Icon(svgEmbedFileMap.GetMustCallback("Bookmark.svg", svgCallback))
	SvgIconCircledChevronDown      = Svg2Icon(svgEmbedFileMap.GetMustCallback("CircledChevronDown.svg", svgCallback))
	SvgIconConvertToNonContainer   = Svg2Icon(svgEmbedFileMap.GetMustCallback("ConvertToNonContainer.svg", svgCallback))
	SvgIconConvertToContainer      = Svg2Icon(svgEmbedFileMap.GetMustCallback("ConvertTo_Container.svg", svgCallback))
	SvgIconLightning               = Svg2Icon(svgEmbedFileMap.GetMustCallback("Lightning.svg", svgCallback))
	SvgIconLogout                  = Svg2Icon(svgEmbedFileMap.GetMustCallback("Logout.svg", svgCallback))
	SvgIconMoon                    = Svg2Icon(svgEmbedFileMap.GetMustCallback("Moon.svg", svgCallback))
	SvgIconSun                     = Svg2Icon(svgEmbedFileMap.GetMustCallback("Sun.svg", svgCallback))
	SvgIconWarning                 = Svg2Icon(svgEmbedFileMap.GetMustCallback("Warning.svg", svgCallback))
	SvgIconAddComment              = Svg2Icon(svgEmbedFileMap.GetMustCallback("add_comment.svg", svgCallback))
	SvgIconAttributes              = Svg2Icon(svgEmbedFileMap.GetMustCallback("attributes.svg", svgCallback))
	SvgIconBack                    = Svg2Icon(svgEmbedFileMap.GetMustCallback("back.svg", svgCallback))
	SvgIconBodyType                = Svg2Icon(svgEmbedFileMap.GetMustCallback("body_type.svg", svgCallback))
	SvgIconBrokenImage             = Svg2Icon(svgEmbedFileMap.GetMustCallback("broken_image.svg", svgCallback))
	SvgIconCalculator              = Svg2Icon(svgEmbedFileMap.GetMustCallback("calculator.svg", svgCallback))
	SvgIconCheckmark               = Svg2Icon(svgEmbedFileMap.GetMustCallback("checkmark.svg", svgCallback))
	SvgIconChevronRight            = Svg2Icon(svgEmbedFileMap.GetMustCallback("chevron_right.svg", svgCallback))
	SvgIconCircledAdd              = Svg2Icon(svgEmbedFileMap.GetMustCallback("circled_add.svg", svgCallback))
	SvgIconCircledChevronRight     = Svg2Icon(svgEmbedFileMap.GetMustCallback("circled_chevron_right.svg", svgCallback))
	SvgIconCircledExclamation      = Svg2Icon(svgEmbedFileMap.GetMustCallback("circled_exclamation.svg", svgCallback))
	SvgIconCircledQuestion         = Svg2Icon(svgEmbedFileMap.GetMustCallback("circled_question.svg", svgCallback))
	SvgIconCircledVerticalEllipsis = Svg2Icon(svgEmbedFileMap.GetMustCallback("circled_vertical_ellipsis.svg", svgCallback))
	SvgIconCircledX                = Svg2Icon(svgEmbedFileMap.GetMustCallback("circled_x.svg", svgCallback))
	SvgIconClone                   = Svg2Icon(svgEmbedFileMap.GetMustCallback("clone.svg", svgCallback))
	SvgIconClosedFolder            = Svg2Icon(svgEmbedFileMap.GetMustCallback("closed_folder.svg", svgCallback))
	SvgIconCoins                   = Svg2Icon(svgEmbedFileMap.GetMustCallback("coins.svg", svgCallback))
	SvgIconContentPasteTwotone     = Svg2Icon(svgEmbedFileMap.GetMustCallback("content_paste_twotone.svg", svgCallback))
	SvgIconCopy                    = Svg2Icon(svgEmbedFileMap.GetMustCallback("copy.svg", svgCallback))
	SvgIconDash                    = Svg2Icon(svgEmbedFileMap.GetMustCallback("dash.svg", svgCallback))
	SvgIconDatabase                = Svg2Icon(svgEmbedFileMap.GetMustCallback("database.svg", svgCallback))
	SvgIconDocument                = Svg2Icon(svgEmbedFileMap.GetMustCallback("document.svg", svgCallback))
	SvgIconDownToBracket           = Svg2Icon(svgEmbedFileMap.GetMustCallback("down_to_bracket.svg", svgCallback))
	SvgIconDownload                = Svg2Icon(svgEmbedFileMap.GetMustCallback("download.svg", svgCallback))
	SvgIconDuplicate               = Svg2Icon(svgEmbedFileMap.GetMustCallback("duplicate.svg", svgCallback))
	SvgIconEdit                    = Svg2Icon(svgEmbedFileMap.GetMustCallback("edit.svg", svgCallback))
	SvgIconFirst                   = Svg2Icon(svgEmbedFileMap.GetMustCallback("first.svg", svgCallback))
	SvgIconFirstAidKit             = Svg2Icon(svgEmbedFileMap.GetMustCallback("first_aid_kit.svg", svgCallback))
	SvgIconFolderSpecialTwotone    = Svg2Icon(svgEmbedFileMap.GetMustCallback("folder_special_twotone.svg", svgCallback))
	SvgIconForward                 = Svg2Icon(svgEmbedFileMap.GetMustCallback("forward.svg", svgCallback))
	SvgIconGcsCampaign             = Svg2Icon(svgEmbedFileMap.GetMustCallback("gcs_campaign.svg", svgCallback))
	SvgIconGcsEquipment            = Svg2Icon(svgEmbedFileMap.GetMustCallback("gcs_equipment.svg", svgCallback))
	SvgIconGcsEquipmentModifiers   = Svg2Icon(svgEmbedFileMap.GetMustCallback("gcs_equipment_modifiers.svg", svgCallback))
	SvgIconGcsLoot                 = Svg2Icon(svgEmbedFileMap.GetMustCallback("gcs_loot.svg", svgCallback))
	SvgIconGcsNotes                = Svg2Icon(svgEmbedFileMap.GetMustCallback("gcs_notes.svg", svgCallback))
	SvgIconGcsSheet                = Svg2Icon(svgEmbedFileMap.GetMustCallback("gcs_sheet.svg", svgCallback))
	SvgIconGcsSkills               = Svg2Icon(svgEmbedFileMap.GetMustCallback("gcs_skills.svg", svgCallback))
	SvgIconGcsSpells               = Svg2Icon(svgEmbedFileMap.GetMustCallback("gcs_spells.svg", svgCallback))
	SvgIconGcsTemplate             = Svg2Icon(svgEmbedFileMap.GetMustCallback("gcs_template.svg", svgCallback))
	SvgIconGcsTraitModifiers       = Svg2Icon(svgEmbedFileMap.GetMustCallback("gcs_trait_modifiers.svg", svgCallback))
	SvgIconGcsTraits               = Svg2Icon(svgEmbedFileMap.GetMustCallback("gcs_traits.svg", svgCallback))
	SvgIconGears                   = Svg2Icon(svgEmbedFileMap.GetMustCallback("gears.svg", svgCallback))
	SvgIconGenericFile             = Svg2Icon(svgEmbedFileMap.GetMustCallback("generic_file.svg", svgCallback))
	SvgIconGrip                    = Svg2Icon(svgEmbedFileMap.GetMustCallback("grip.svg", svgCallback))
	SvgIconHelp                    = Svg2Icon(svgEmbedFileMap.GetMustCallback("help.svg", svgCallback))
	SvgIconHierarchy               = Svg2Icon(svgEmbedFileMap.GetMustCallback("hierarchy.svg", svgCallback))
	SvgIconImageFile               = Svg2Icon(svgEmbedFileMap.GetMustCallback("image_file.svg", svgCallback))
	SvgIconInfo                    = Svg2Icon(svgEmbedFileMap.GetMustCallback("info.svg", svgCallback))
	SvgIconLast                    = Svg2Icon(svgEmbedFileMap.GetMustCallback("last.svg", svgCallback))
	SvgIconLink                    = Svg2Icon(svgEmbedFileMap.GetMustCallback("link.svg", svgCallback))
	SvgIconMarkdownFile            = Svg2Icon(svgEmbedFileMap.GetMustCallback("markdown_file.svg", svgCallback))
	SvgIconMediation               = Svg2Icon(svgEmbedFileMap.GetMustCallback("mediation.svg", svgCallback))
	SvgIconMeleeWeapon             = Svg2Icon(svgEmbedFileMap.GetMustCallback("melee_weapon.svg", svgCallback))
	SvgIconMenu                    = Svg2Icon(svgEmbedFileMap.GetMustCallback("menu.svg", svgCallback))
	SvgIconNaming                  = Svg2Icon(svgEmbedFileMap.GetMustCallback("naming.svg", svgCallback))
	SvgIconNewFolder               = Svg2Icon(svgEmbedFileMap.GetMustCallback("new_folder.svg", svgCallback))
	SvgIconNext                    = Svg2Icon(svgEmbedFileMap.GetMustCallback("next.svg", svgCallback))
	SvgIconNot                     = Svg2Icon(svgEmbedFileMap.GetMustCallback("not.svg", svgCallback))
	SvgIconNotesCollapse           = Svg2Icon(svgEmbedFileMap.GetMustCallback("notes-collapse.svg", svgCallback))
	SvgIconNotesExpand             = Svg2Icon(svgEmbedFileMap.GetMustCallback("notes-expand.svg", svgCallback))
	SvgIconOpenFolder              = Svg2Icon(svgEmbedFileMap.GetMustCallback("open_folder.svg", svgCallback))
	SvgIconPdfFile                 = Svg2Icon(svgEmbedFileMap.GetMustCallback("pdf_file.svg", svgCallback))
	SvgIconPolylineRound           = Svg2Icon(svgEmbedFileMap.GetMustCallback("polyline_round.svg", svgCallback))
	SvgIconPrevious                = Svg2Icon(svgEmbedFileMap.GetMustCallback("previous.svg", svgCallback))
	SvgIconRandomize               = Svg2Icon(svgEmbedFileMap.GetMustCallback("randomize.svg", svgCallback))
	SvgIconRangedWeapon            = Svg2Icon(svgEmbedFileMap.GetMustCallback("ranged_weapon.svg", svgCallback))
	SvgIconReleaseNotes            = Svg2Icon(svgEmbedFileMap.GetMustCallback("release_notes.svg", svgCallback))
	SvgIconReset                   = Svg2Icon(svgEmbedFileMap.GetMustCallback("reset.svg", svgCallback))
	SvgIconSaveContent             = Svg2Icon(svgEmbedFileMap.GetMustCallback("save_content.svg", svgCallback))
	SvgIconSettings                = Svg2Icon(svgEmbedFileMap.GetMustCallback("settings.svg", svgCallback))
	SvgIconSideBar                 = Svg2Icon(svgEmbedFileMap.GetMustCallback("side_bar.svg", svgCallback))
	SvgIconSignPost                = Svg2Icon(svgEmbedFileMap.GetMustCallback("sign_post.svg", svgCallback))
	SvgIconSizeToFit               = Svg2Icon(svgEmbedFileMap.GetMustCallback("size_to_fit.svg", svgCallback))
	SvgIconSortAscending           = Svg2Icon(svgEmbedFileMap.GetMustCallback("sort_ascending.svg", svgCallback))
	SvgIconSortDescending          = Svg2Icon(svgEmbedFileMap.GetMustCallback("sort_descending.svg", svgCallback))
	SvgIconStack                   = Svg2Icon(svgEmbedFileMap.GetMustCallback("stack.svg", svgCallback))
	SvgIconStamper                 = Svg2Icon(svgEmbedFileMap.GetMustCallback("stamper.svg", svgCallback))
	SvgIconStar                    = Svg2Icon(svgEmbedFileMap.GetMustCallback("star.svg", svgCallback))
	SvgIconTrash                   = Svg2Icon(svgEmbedFileMap.GetMustCallback("trash.svg", svgCallback))
	SvgIconTriangleExclamation     = Svg2Icon(svgEmbedFileMap.GetMustCallback("triangle_exclamation.svg", svgCallback))
	SvgIconWeight                  = Svg2Icon(svgEmbedFileMap.GetMustCallback("weight.svg", svgCallback))
	SvgIconWindowMaximize          = Svg2Icon(svgEmbedFileMap.GetMustCallback("window_maximize.svg", svgCallback))
	SvgIconWindowRestore           = Svg2Icon(svgEmbedFileMap.GetMustCallback("window_restore.svg", svgCallback))
)
