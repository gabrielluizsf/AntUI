package dev.antui;

import android.content.ContentProvider;
import android.content.ContentValues;
import android.database.Cursor;
import android.database.MatrixCursor;
import android.net.Uri;
import android.os.ParcelFileDescriptor;
import android.provider.OpenableColumns;
import android.webkit.MimeTypeMap;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;

/**
 * Hands a file to another app.
 *
 * <p>Since Android 7 no app may give another a path: a file:// address in an
 * intent throws, and the only way to share anything is a content:// address
 * from a provider. The one everybody uses is androidx.core's FileProvider,
 * which means a dependency and an XML resource; this is the same idea in
 * eighty lines, serving one directory and refusing everything else.
 *
 * <p>That one directory is the point. A provider that serves whatever path it
 * is handed is a provider that hands out the app's private data to whoever
 * asks — so this one resolves every request under a single folder inside the
 * cache and checks the canonical path, which is what stops "../../databases"
 * from being a valid request.
 */
public final class AntuiFileProvider extends ContentProvider {

    /** Appended to the app's package to make the authority. */
    public static final String SUFFIX = ".antuifiles";

    /** The only directory this provider will serve, inside the cache. */
    public static final String DIR = "antui-share";

    private File root;

    @Override
    public boolean onCreate() {
        root = new File(getContext().getCacheDir(), DIR);
        root.mkdirs();
        return true;
    }

    private File resolve(Uri uri) throws FileNotFoundException {
        String path = uri.getPath();
        if (path == null || root == null) {
            throw new FileNotFoundException(String.valueOf(uri));
        }
        File file = new File(root, path);
        try {
            String inside = root.getCanonicalPath() + File.separator;
            if (!file.getCanonicalPath().startsWith(inside)) {
                // A path that climbs out. It is not an error worth
                // explaining to the caller: as far as it is concerned there
                // is no such file.
                throw new FileNotFoundException(String.valueOf(uri));
            }
        } catch (IOException e) {
            throw new FileNotFoundException(String.valueOf(uri));
        }
        if (!file.isFile()) {
            throw new FileNotFoundException(String.valueOf(uri));
        }
        return file;
    }

    @Override
    public ParcelFileDescriptor openFile(Uri uri, String mode) throws FileNotFoundException {
        // Read only, whatever was asked for. Nothing this app shares is
        // something another app should be writing.
        return ParcelFileDescriptor.open(resolve(uri), ParcelFileDescriptor.MODE_READ_ONLY);
    }

    @Override
    public Cursor query(Uri uri, String[] projection, String selection,
                        String[] selectionArgs, String sortOrder) {
        File file;
        try {
            file = resolve(uri);
        } catch (FileNotFoundException e) {
            return null;
        }
        String[] columns = projection;
        if (columns == null) {
            columns = new String[]{OpenableColumns.DISPLAY_NAME, OpenableColumns.SIZE};
        }
        Object[] row = new Object[columns.length];
        for (int i = 0; i < columns.length; i++) {
            if (OpenableColumns.DISPLAY_NAME.equals(columns[i])) {
                row[i] = file.getName();
            } else if (OpenableColumns.SIZE.equals(columns[i])) {
                row[i] = Long.valueOf(file.length());
            }
        }
        MatrixCursor cursor = new MatrixCursor(columns, 1);
        cursor.addRow(row);
        return cursor;
    }

    @Override
    public String getType(Uri uri) {
        String name = uri.getLastPathSegment();
        if (name == null) {
            return null;
        }
        int dot = name.lastIndexOf('.');
        if (dot < 0 || dot == name.length() - 1) {
            return "application/octet-stream";
        }
        String extension = name.substring(dot + 1).toLowerCase();
        String mime = MimeTypeMap.getSingleton().getMimeTypeFromExtension(extension);
        return mime != null ? mime : "application/octet-stream";
    }

    // Nothing else is offered. A provider that only hands over files does not
    // need to pretend to be a database.

    @Override
    public Uri insert(Uri uri, ContentValues values) {
        throw new UnsupportedOperationException("antui: this provider only reads");
    }

    @Override
    public int delete(Uri uri, String selection, String[] selectionArgs) {
        throw new UnsupportedOperationException("antui: this provider only reads");
    }

    @Override
    public int update(Uri uri, ContentValues values, String selection, String[] args) {
        throw new UnsupportedOperationException("antui: this provider only reads");
    }
}
