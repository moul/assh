//Copyright 2017 GoGraphviz Authors
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http)://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package gographviz

import "fmt"

// Attr is an attribute key
type Attr string

// NewAttr creates a new attribute key by checking whether it is a valid key
func NewAttr(key string) (Attr, error) {
	a, ok := validAttrs[key]
	if !ok {
		return Attr(""), fmt.Errorf("%s is not a valid attribute", key)
	}
	return a, nil
}

const (
	// Damping http://www.graphviz.org/content/attrs#dDamping
	Damping Attr = "Damping"
	// K http://www.graphviz.org/content/attrs#dK
	K Attr = "K"
	// URL http://www.graphviz.org/content/attrs#dURL
	URL Attr = "URL"
	// Background http://www.graphviz.org/content/attrs#d_background
	Background Attr = "_background"
	// Area http://www.graphviz.org/content/attrs#darea
	Area Attr = "area"
	// ArrowHead http://www.graphviz.org/content/attrs#darrowhead
	ArrowHead Attr = "arrowhead"
	// ArrowSize http://www.graphviz.org/content/attrs#darrowsize
	ArrowSize Attr = "arrowsize"
	// ArrowTail http://www.graphviz.org/content/attrs#darrowtail
	ArrowTail Attr = "arrowtail"
	// BB http://www.graphviz.org/content/attrs#dbb
	BB Attr = "bb"
	// BgColor http://www.graphviz.org/content/attrs#dbgcolor
	BgColor Attr = "bgcolor"
	// Center http://www.graphviz.org/content/attrs#dcenter
	Center Attr = "center"
	// Charset http://www.graphviz.org/content/attrs#dcharset
	Charset Attr = "charset"
	// ClusterRank http://www.graphviz.org/content/attrs#dclusterrank
	ClusterRank Attr = "clusterrank"
	// Color http://www.graphviz.org/content/attrs#dcolor
	Color Attr = "color"
	// ColorScheme http://www.graphviz.org/content/attrs#dcolorscheme
	ColorScheme Attr = "colorscheme"
	// Comment http://www.graphviz.org/content/attrs#dcomment
	Comment Attr = "comment"
	// Compound http://www.graphviz.org/content/attrs#dcompound
	Compound Attr = "compound"
	// Concentrate http://www.graphviz.org/content/attrs#dconcentrate
	Concentrate Attr = "concentrate"
	// Constraint http://www.graphviz.org/content/attrs#dconstraint
	Constraint Attr = "constraint"
	// Decorate http://www.graphviz.org/content/attrs#ddecorate
	Decorate Attr = "decorate"
	// DefaultDist http://www.graphviz.org/content/attrs#ddefaultdist
	DefaultDist Attr = "defaultdist"
	// Dim http://www.graphviz.org/content/attrs#ddim
	Dim Attr = "dim"
	// Dimen http://www.graphviz.org/content/attrs#ddimen
	Dimen Attr = "dimen"
	// Dir http://www.graphviz.org/content/attrs#ddir
	Dir Attr = "dir"
	// DirEdgeConstraints http://www.graphviz.org/content/attrs#ddir
	DirEdgeConstraints Attr = "diredgeconstraints"
	// Distortion http://www.graphviz.org/content/attrs#ddistortion
	Distortion Attr = "distortion"
	// DPI http://www.graphviz.org/content/attrs#ddpi
	DPI Attr = "dpi"
	// EdgeURL http://www.graphviz.org/content/attrs#d:edgeURL
	EdgeURL Attr = "edgeURL"
	// EdgeHREF http://www.graphviz.org/content/attrs#d:edgehref
	EdgeHREF Attr = "edgehref"
	// EdgeTarget http://www.graphviz.org/content/attrs#d:edgetarget
	EdgeTarget Attr = "edgetarget"
	// EdgeTooltip http://www.graphviz.org/content/attrs#d:edgetooltip
	EdgeTooltip Attr = "edgetooltip"
	// Epsilon http://www.graphviz.org/content/attrs#d:epsilon
	Epsilon Attr = "epsilon"
	// ESep http://www.graphviz.org/content/attrs#d:epsilon
	ESep Attr = "esep"
	// FillColor http://www.graphviz.org/content/attrs#dfillcolor
	FillColor Attr = "fillcolor"
	// FixedSize http://www.graphviz.org/content/attrs#dfixedsize
	FixedSize Attr = "fixedsize"
	// FontColor http://www.graphviz.org/content/attrs#dfontcolor
	FontColor Attr = "fontcolor"
	// FontName http://www.graphviz.org/content/attrs#dfontname
	FontName Attr = "fontname"
	// FontNames http://www.graphviz.org/content/attrs#dfontnames
	FontNames Attr = "fontnames"
	// FontPath http://www.graphviz.org/content/attrs#dfontpath
	FontPath Attr = "fontpath"
	// FontSize http://www.graphviz.org/content/attrs#dfontsize
	FontSize Attr = "fontsize"
	// ForceLabels http://www.graphviz.org/content/attrs#dforcelabels
	ForceLabels Attr = "forcelabels"
	// GradientAngle http://www.graphviz.org/content/attrs#dgradientangle
	GradientAngle Attr = "gradientangle"
	// Group http://www.graphviz.org/content/attrs#dgroup
	Group Attr = "group"
	// HeadURL http://www.graphviz.org/content/attrs#dheadURL
	HeadURL Attr = "headURL"
	// HeadLP http://www.graphviz.org/content/attrs#dhead_lp
	HeadLP Attr = "head_lp"
	// HeadClip http://www.graphviz.org/content/attrs#dheadclip
	HeadClip Attr = "headclip"
	// HeadHREF http://www.graphviz.org/content/attrs#dheadhref
	HeadHREF Attr = "headhref"
	// HeadLabel http://www.graphviz.org/content/attrs#dheadlabel
	HeadLabel Attr = "headlabel"
	// HeadPort http://www.graphviz.org/content/attrs#dheadport
	HeadPort Attr = "headport"
	// HeadTarget http://www.graphviz.org/content/attrs#dheadtarget
	HeadTarget Attr = "headtarget"
	// HeadTooltip http://www.graphviz.org/content/attrs#dheadtooltip
	HeadTooltip Attr = "headtooltip"
	// Height http://www.graphviz.org/content/attrs#dheight
	Height Attr = "height"
	// HREF http://www.graphviz.org/content/attrs#dhref
	HREF Attr = "href"
	// ID http://www.graphviz.org/content/attrs#did
	ID Attr = "id"
	// Image http://www.graphviz.org/content/attrs#dimage
	Image Attr = "image"
	// ImagePath http://www.graphviz.org/content/attrs#dimagepath
	ImagePath Attr = "imagepath"
	// ImageScale http://www.graphviz.org/content/attrs#dimagescale
	ImageScale Attr = "imagescale"
	// InputScale http://www.graphviz.org/content/attrs#dinputscale
	InputScale Attr = "inputscale"
	// Label http://www.graphviz.org/content/attrs#dlabel
	Label Attr = "label"
	// LabelURL http://www.graphviz.org/content/attrs#dlabelURL
	LabelURL Attr = "labelURL"
	// LabelScheme http://www.graphviz.org/content/attrs#dlabel_scheme
	LabelScheme Attr = "label_scheme"
	// LabelAngle http://www.graphviz.org/content/attrs#dlabelangle
	LabelAngle Attr = "labelangle"
	// LabelDistance http://www.graphviz.org/content/attrs#dlabeldistance
	LabelDistance Attr = "labeldistance"
	// LabelFloat http://www.graphviz.org/content/attrs#dlabelfloat
	LabelFloat Attr = "labelfloat"
	// LabelFontColor http://www.graphviz.org/content/attrs#dlabelfontcolor
	LabelFontColor Attr = "labelfontcolor"
	// LabelFontName http://www.graphviz.org/content/attrs#dlabelfontname
	LabelFontName Attr = "labelfontname"
	// LabelFontSize http://www.graphviz.org/content/attrs#dlabelfontsize
	LabelFontSize Attr = "labelfontsize"
	// LabelHREF http://www.graphviz.org/content/attrs#dlabelhref
	LabelHREF Attr = "labelhref"
	// LabelJust http://www.graphviz.org/content/attrs#dlabeljust
	LabelJust Attr = "labeljust"
	// LabelLOC http://www.graphviz.org/content/attrs#dlabelloc
	LabelLOC Attr = "labelloc"
	// LabelTarget http://www.graphviz.org/content/attrs#dlabeltarget
	LabelTarget Attr = "labeltarget"
	// LabelTooltip http://www.graphviz.org/content/attrs#dlabeltooltip
	LabelTooltip Attr = "labeltooltip"
	// Landscape http://www.graphviz.org/content/attrs#dlandscape
	Landscape Attr = "landscape"
	// Layer http://www.graphviz.org/content/attrs#dlayer
	Layer Attr = "layer"
	// LayerListSep http://www.graphviz.org/content/attrs#dlayerlistsep
	LayerListSep Attr = "layerlistsep"
	// Layers http://www.graphviz.org/content/attrs#dlayers
	Layers Attr = "layers"
	// LayerSelect http://www.graphviz.org/content/attrs#dlayerselect
	LayerSelect Attr = "layerselect"
	// LayerSep http://www.graphviz.org/content/attrs#dlayersep
	LayerSep Attr = "layersep"
	// Layout http://www.graphviz.org/content/attrs#dlayout
	Layout Attr = "layout"
	// Len http://www.graphviz.org/content/attrs#dlen
	Len Attr = "len"
	// Levels http://www.graphviz.org/content/attrs#dlevels
	Levels Attr = "levels"
	// LevelsGap http://www.graphviz.org/content/attrs#dlevelsgap
	LevelsGap Attr = "levelsgap"
	// LHead http://www.graphviz.org/content/attrs#dlhead
	LHead Attr = "lhead"
	// LHeight http://www.graphviz.org/content/attrs#dlheight
	LHeight Attr = "lheight"
	// LP http://www.graphviz.org/content/attrs#dlp
	LP Attr = "lp"
	// LTail http://www.graphviz.org/content/attrs#dltail
	LTail Attr = "ltail"
	// LWidth http://www.graphviz.org/content/attrs#dlwidth
	LWidth Attr = "lwidth"
	// Margin http://www.graphviz.org/content/attrs#dmargin
	Margin Attr = "margin"
	// MaxIter http://www.graphviz.org/content/attrs#dmaxiter
	MaxIter Attr = "maxiter"
	// MCLimit http://www.graphviz.org/content/attrs#dmclimit
	MCLimit Attr = "mclimit"
	// MinDist http://www.graphviz.org/content/attrs#dmindist
	MinDist Attr = "mindist"
	// MinLen http://www.graphviz.org/content/attrs#dmindist
	MinLen Attr = "minlen"
	// Mode http://www.graphviz.org/content/attrs#dmode
	Mode Attr = "mode"
	// Model http://www.graphviz.org/content/attrs#dmodel
	Model Attr = "model"
	// Mosek http://www.graphviz.org/content/attrs#dmosek
	Mosek Attr = "mosek"
	// NewRank http://www.graphviz.org/content/attrs#dnewrank
	NewRank Attr = "newrank"
	// NodeSep http://www.graphviz.org/content/attrs#dnodesep
	NodeSep Attr = "nodesep"
	// NoJustify http://www.graphviz.org/content/attrs#dnojustify
	NoJustify Attr = "nojustify"
	// Normalize http://www.graphviz.org/content/attrs#dnormalize
	Normalize Attr = "normalize"
	// NoTranslate http://www.graphviz.org/content/attrs#dnotranslate
	NoTranslate Attr = "notranslate"
	// NSLimit http://www.graphviz.org/content/attrs#dnslimit
	NSLimit Attr = "nslimit"
	// NSLimit1 http://www.graphviz.org/content/attrs#dnslimit1
	NSLimit1 Attr = "nslimit1"
	// Ordering http://www.graphviz.org/content/attrs#dnslimit1
	Ordering Attr = "ordering"
	// Orientation http://www.graphviz.org/content/attrs#dorientation
	Orientation Attr = "orientation"
	// OutputOrder http://www.graphviz.org/content/attrs#doutputorder
	OutputOrder Attr = "outputorder"
	// Overlap http://www.graphviz.org/content/attrs#doverlap
	Overlap Attr = "overlap"
	// OverlapScaling http://www.graphviz.org/content/attrs#doverlap_scaling
	OverlapScaling Attr = "overlap_scaling"
	// OverlapShrink http://www.graphviz.org/content/attrs#doverlap_shrink
	OverlapShrink Attr = "overlap_shrink"
	// Pack http://www.graphviz.org/content/attrs#dpack
	Pack Attr = "pack"
	// PackMode http://www.graphviz.org/content/attrs#dpackmode
	PackMode Attr = "packmode"
	// Pad http://www.graphviz.org/content/attrs#dpad
	Pad Attr = "pad"
	// Page http://www.graphviz.org/content/attrs#dpage
	Page Attr = "page"
	// PageDir http://www.graphviz.org/content/attrs#dpagedir
	PageDir Attr = "pagedir"
	// PenColor http://www.graphviz.org/content/attrs#dpencolor
	PenColor Attr = "pencolor"
	// PenWidth http://www.graphviz.org/content/attrs#dpenwidth
	PenWidth Attr = "penwidth"
	// Peripheries http://www.graphviz.org/content/attrs#dperipheries
	Peripheries Attr = "peripheries"
	// Pin http://www.graphviz.org/content/attrs#dperipheries
	Pin Attr = "pin"
	// Pos http://www.graphviz.org/content/attrs#dpos
	Pos Attr = "pos"
	// QuadTree http://www.graphviz.org/content/attrs#dquadtree
	QuadTree Attr = "quadtree"
	// Quantum http://www.graphviz.org/content/attrs#dquantum
	Quantum Attr = "quantum"
	// Rank http://www.graphviz.org/content/attrs#drank
	Rank Attr = "rank"
	// RankDir http://www.graphviz.org/content/attrs#drankdir
	RankDir Attr = "rankdir"
	// RankSep http://www.graphviz.org/content/attrs#dranksep
	RankSep Attr = "ranksep"
	// Ratio http://www.graphviz.org/content/attrs#dratio
	Ratio Attr = "ratio"
	// Rects http://www.graphviz.org/content/attrs#drects
	Rects Attr = "rects"
	// Regular http://www.graphviz.org/content/attrs#dregular
	Regular Attr = "regular"
	// ReMinCross http://www.graphviz.org/content/attrs#dremincross
	ReMinCross Attr = "remincross"
	// RepulsiveForce http://www.graphviz.org/content/attrs#drepulsiveforce
	RepulsiveForce Attr = "repulsiveforce"
	// Resolution http://www.graphviz.org/content/attrs#dresolution
	Resolution Attr = "resolution"
	// Root http://www.graphviz.org/content/attrs#droot
	Root Attr = "root"
	// Rotate http://www.graphviz.org/content/attrs#drotate
	Rotate Attr = "rotate"
	// Rotation http://www.graphviz.org/content/attrs#drotation
	Rotation Attr = "rotation"
	// SameHead http://www.graphviz.org/content/attrs#dsamehead
	SameHead Attr = "samehead"
	// SameTail http://www.graphviz.org/content/attrs#dsametail
	SameTail Attr = "sametail"
	// SamplePoints http://www.graphviz.org/content/attrs#dsamplepoints
	SamplePoints Attr = "samplepoints"
	// Scale http://www.graphviz.org/content/attrs#dscale
	Scale Attr = "scale"
	// SearchSize http://www.graphviz.org/content/attrs#dsearchsize
	SearchSize Attr = "searchsize"
	// Sep http://www.graphviz.org/content/attrs#dsep
	Sep Attr = "sep"
	// Shape http://www.graphviz.org/content/attrs#dshape
	Shape Attr = "shape"
	// ShapeFile http://www.graphviz.org/content/attrs#dshapefile
	ShapeFile Attr = "shapefile"
	// ShowBoxes http://www.graphviz.org/content/attrs#dshowboxes
	ShowBoxes Attr = "showboxes"
	// Sides http://www.graphviz.org/content/attrs#dsides
	Sides Attr = "sides"
	// Size http://www.graphviz.org/content/attrs#dsize
	Size Attr = "size"
	// Skew http://www.graphviz.org/content/attrs#dskew
	Skew Attr = "skew"
	// Smoothing http://www.graphviz.org/content/attrs#dsmoothing
	Smoothing Attr = "smoothing"
	// SortV http://www.graphviz.org/content/attrs#dsortv
	SortV Attr = "sortv"
	// Splines http://www.graphviz.org/content/attrs#dsplines
	Splines Attr = "splines"
	// Start http://www.graphviz.org/content/attrs#dstart
	Start Attr = "start"
	// Style http://www.graphviz.org/content/attrs#dstyle
	Style Attr = "style"
	// StyleSheet http://www.graphviz.org/content/attrs#dstylesheet
	StyleSheet Attr = "stylesheet"
	// TailURL http://www.graphviz.org/content/attrs#dtailURL
	TailURL Attr = "tailURL"
	// TailLP http://www.graphviz.org/content/attrs#dtail_lp
	TailLP Attr = "tail_lp"
	// TailClip http://www.graphviz.org/content/attrs#dtailclip
	TailClip Attr = "tailclip"
	// TailHREF http://www.graphviz.org/content/attrs#dtailhref
	TailHREF Attr = "tailhref"
	// TailLabel http://www.graphviz.org/content/attrs#dtaillabel
	TailLabel Attr = "taillabel"
	// TailPort http://www.graphviz.org/content/attrs#dtailport
	TailPort Attr = "tailport"
	// TailTarget http://www.graphviz.org/content/attrs#dtailtarget
	TailTarget Attr = "tailtarget"
	// TailTooltip http://www.graphviz.org/content/attrs#dtailtooltip
	TailTooltip Attr = "tailtooltip"
	// Target http://www.graphviz.org/content/attrs#dtarget
	Target Attr = "target"
	// Tooltip http://www.graphviz.org/content/attrs#dtooltip
	Tooltip Attr = "tooltip"
	// TrueColor http://www.graphviz.org/content/attrs#dtooltip
	TrueColor Attr = "truecolor"
	// Vertices http://www.graphviz.org/content/attrs#dvertices
	Vertices Attr = "vertices"
	// ViewPort http://www.graphviz.org/content/attrs#dviewport
	ViewPort Attr = "viewport"
	// VoroMargin http://www.graphviz.org/content/attrs#dvoro_margin
	VoroMargin Attr = "voro_margin"
	// Weight http://www.graphviz.org/content/attrs#dweight
	Weight Attr = "weight"
	// Width http://www.graphviz.org/content/attrs#dwidth
	Width Attr = "width"
	// XDotVersion http://www.graphviz.org/content/attrs#dxdotversion
	XDotVersion Attr = "xdotversion"
	// XLabel http://www.graphviz.org/content/attrs#dxlabel
	XLabel Attr = "xlabel"
	// XLP http://www.graphviz.org/content/attrs#dxlp
	XLP Attr = "xlp"
	// Z http://www.graphviz.org/content/attrs#dz
	Z Attr = "z"

	// MinCross is not in the documentation, but found in the Ped_Lion_Share (lion_share.gv.txt) example
	MinCross Attr = "mincross"
	// SSize is not in the documentation, but found in the siblings.gv.txt example
	SSize Attr = "ssize"
	// Outline is not in the documentation, but found in the siblings.gv.txt example
	Outline Attr = "outline"
	// F is not in the documentation, but found in the transparency.gv.txt example
	F Attr = "f"
)

var validAttrs = map[string]Attr{
	string(Damping):            Damping,
	string(K):                  K,
	string(URL):                URL,
	string(Background):         Background,
	string(Area):               Area,
	string(ArrowHead):          ArrowHead,
	string(ArrowSize):          ArrowSize,
	string(ArrowTail):          ArrowTail,
	string(BB):                 BB,
	string(BgColor):            BgColor,
	string(Center):             Center,
	string(Charset):            Charset,
	string(ClusterRank):        ClusterRank,
	string(Color):              Color,
	string(ColorScheme):        ColorScheme,
	string(Comment):            Comment,
	string(Compound):           Compound,
	string(Concentrate):        Concentrate,
	string(Constraint):         Constraint,
	string(Decorate):           Decorate,
	string(DefaultDist):        DefaultDist,
	string(Dim):                Dim,
	string(Dimen):              Dimen,
	string(Dir):                Dir,
	string(DirEdgeConstraints): DirEdgeConstraints,
	string(Distortion):         Distortion,
	string(DPI):                DPI,
	string(EdgeURL):            EdgeURL,
	string(EdgeHREF):           EdgeHREF,
	string(EdgeTarget):         EdgeTarget,
	string(EdgeTooltip):        EdgeTooltip,
	string(Epsilon):            Epsilon,
	string(ESep):               ESep,
	string(FillColor):          FillColor,
	string(FixedSize):          FixedSize,
	string(FontColor):          FontColor,
	string(FontName):           FontName,
	string(FontNames):          FontNames,
	string(FontPath):           FontPath,
	string(FontSize):           FontSize,
	string(ForceLabels):        ForceLabels,
	string(GradientAngle):      GradientAngle,
	string(Group):              Group,
	string(HeadURL):            HeadURL,
	string(HeadLP):             HeadLP,
	string(HeadClip):           HeadClip,
	string(HeadHREF):           HeadHREF,
	string(HeadLabel):          HeadLabel,
	string(HeadPort):           HeadPort,
	string(HeadTarget):         HeadTarget,
	string(HeadTooltip):        HeadTooltip,
	string(Height):             Height,
	string(HREF):               HREF,
	string(ID):                 ID,
	string(Image):              Image,
	string(ImagePath):          ImagePath,
	string(ImageScale):         ImageScale,
	string(InputScale):         InputScale,
	string(Label):              Label,
	string(LabelURL):           LabelURL,
	string(LabelScheme):        LabelScheme,
	string(LabelAngle):         LabelAngle,
	string(LabelDistance):      LabelDistance,
	string(LabelFloat):         LabelFloat,
	string(LabelFontColor):     LabelFontColor,
	string(LabelFontName):      LabelFontName,
	string(LabelFontSize):      LabelFontSize,
	string(LabelHREF):          LabelHREF,
	string(LabelJust):          LabelJust,
	string(LabelLOC):           LabelLOC,
	string(LabelTarget):        LabelTarget,
	string(LabelTooltip):       LabelTooltip,
	string(Landscape):          Landscape,
	string(Layer):              Layer,
	string(LayerListSep):       LayerListSep,
	string(Layers):             Layers,
	string(LayerSelect):        LayerSelect,
	string(LayerSep):           LayerSep,
	string(Layout):             Layout,
	string(Len):                Len,
	string(Levels):             Levels,
	string(LevelsGap):          LevelsGap,
	string(LHead):              LHead,
	string(LHeight):            LHeight,
	string(LP):                 LP,
	string(LTail):              LTail,
	string(LWidth):             LWidth,
	string(Margin):             Margin,
	string(MaxIter):            MaxIter,
	string(MCLimit):            MCLimit,
	string(MinDist):            MinDist,
	string(MinLen):             MinLen,
	string(Mode):               Mode,
	string(Model):              Model,
	string(Mosek):              Mosek,
	string(NewRank):            NewRank,
	string(NodeSep):            NodeSep,
	string(NoJustify):          NoJustify,
	string(Normalize):          Normalize,
	string(NoTranslate):        NoTranslate,
	string(NSLimit):            NSLimit,
	string(NSLimit1):           NSLimit1,
	string(Ordering):           Ordering,
	string(Orientation):        Orientation,
	string(OutputOrder):        OutputOrder,
	string(Overlap):            Overlap,
	string(OverlapScaling):     OverlapScaling,
	string(OverlapShrink):      OverlapShrink,
	string(Pack):               Pack,
	string(PackMode):           PackMode,
	string(Pad):                Pad,
	string(Page):               Page,
	string(PageDir):            PageDir,
	string(PenColor):           PenColor,
	string(PenWidth):           PenWidth,
	string(Peripheries):        Peripheries,
	string(Pin):                Pin,
	string(Pos):                Pos,
	string(QuadTree):           QuadTree,
	string(Quantum):            Quantum,
	string(Rank):               Rank,
	string(RankDir):            RankDir,
	string(RankSep):            RankSep,
	string(Ratio):              Ratio,
	string(Rects):              Rects,
	string(Regular):            Regular,
	string(ReMinCross):         ReMinCross,
	string(RepulsiveForce):     RepulsiveForce,
	string(Resolution):         Resolution,
	string(Root):               Root,
	string(Rotate):             Rotate,
	string(Rotation):           Rotation,
	string(SameHead):           SameHead,
	string(SameTail):           SameTail,
	string(SamplePoints):       SamplePoints,
	string(Scale):              Scale,
	string(SearchSize):         SearchSize,
	string(Sep):                Sep,
	string(Shape):              Shape,
	string(ShapeFile):          ShapeFile,
	string(ShowBoxes):          ShowBoxes,
	string(Sides):              Sides,
	string(Size):               Size,
	string(Skew):               Skew,
	string(Smoothing):          Smoothing,
	string(SortV):              SortV,
	string(Splines):            Splines,
	string(Start):              Start,
	string(Style):              Style,
	string(StyleSheet):         StyleSheet,
	string(TailURL):            TailURL,
	string(TailLP):             TailLP,
	string(TailClip):           TailClip,
	string(TailHREF):           TailHREF,
	string(TailLabel):          TailLabel,
	string(TailPort):           TailPort,
	string(TailTarget):         TailTarget,
	string(TailTooltip):        TailTooltip,
	string(Target):             Target,
	string(Tooltip):            Tooltip,
	string(TrueColor):          TrueColor,
	string(Vertices):           Vertices,
	string(ViewPort):           ViewPort,
	string(VoroMargin):         VoroMargin,
	string(Weight):             Weight,
	string(Width):              Width,
	string(XDotVersion):        XDotVersion,
	string(XLabel):             XLabel,
	string(XLP):                XLP,
	string(Z):                  Z,

	string(MinCross): MinCross,
	string(SSize):    SSize,
	string(Outline):  Outline,
	string(F):        F,
}
