package apk

// What the store needs that no build produces.
//
// [Build] and [Bundle] make the file. The file is perhaps a third of what a
// release is: the rest is a listing, a set of declarations, and a review,
// none of which anything here can make and all of which stop an app going
// out. They are written down here rather than in a wiki because a wiki goes
// stale silently and a constant in the tool people run does not.
//
// [Checklist] is the text; "antuiapk publish" prints it.

// Checklist is what a human has to do that no build can.
const Checklist = `Before the first release

  A developer account, which costs $25 once, and identity verification —
  a legal name, an address and a phone number that are checked. This takes
  days rather than minutes and cannot be started after the app is ready.

  If the account is a personal one rather than an organisation, the app
  must be tested by 12 people for 14 continuous days before it may go to
  production. There is no way round it and no way to shorten it, so it is
  the thing to start first.

The listing, none of which is in the package

  App name, 30 characters. Short description, 80. Full description, 4000.
  Icon, 512x512 PNG, 32-bit, under 1 MB. This is not the icon in the APK:
  the store shows its own and the two are uploaded separately.
  Feature graphic, 1024x500, shown at the top of the listing.
  At least two phone screenshots, between 320 and 3840 pixels on a side.
  Screenshots for tablets, TV or Wear if the app is offered to them.

Declarations, each of which blocks the release until it is answered

  A privacy policy at a URL that works. Every app needs one now, whatever
  permissions it asks for.
  The Data safety form: what the app collects, what it shares, whether it
  is encrypted in transit, whether the user can ask for deletion. It is
  answered by hand and it is checked against what the app does.
  The content rating questionnaire, which issues an IARC rating.
  Target audience and age. An app that says it is for children is held to
  a great deal more.
  Whether the app has ads.
  App access: if anything is behind a login, working credentials for the
  reviewer, or the review fails on a sign-in screen.

Signing, which is not what it looks like

  A new app uses Play App Signing, and that changes what the key is. The
  key made by "antuiapk keygen" and used to sign the bundle is the *upload*
  key: Google strips that signature and re-signs with an app signing key it
  holds and never shows anyone. So losing the upload key is recoverable —
  ask the Console to register a new one — and the app signing key cannot be
  lost, which is the point of the arrangement.

  It also means the certificate a device sees is not the one here. Anything
  that depends on the signature — App Links assetlinks.json, an API key
  restricted by fingerprint, another app checking who signed this one — must
  use the app signing certificate from the Console, not the one from the
  keystore.

Uploading

  By hand, in the Play Console: make a release on a track, drop the .aab in,
  add release notes, roll out.

  antuiapk does not upload. There is an API for it — the Google Play
  Developer API, with a service account from Google Cloud that has to be
  granted access in the Console, and an upload that is an edit opened, a
  bundle posted, a track set and the edit committed — but it is not here,
  because it cannot be tested without a real account and a real app, and an
  untested thing that handles a signing credential and publishes to the
  public is worse than no thing at all.

Once, before the first upload, and never again

  The application id. It cannot be changed, and it cannot be reused by
  anyone including you.

Every time

  The version code goes up. A code that has been uploaded is spent forever,
  including by a rejected upload.
  The symbols go with it: antuiapk puts them in the bundle, and a native
  crash without them is a list of addresses that can never be read.
`
