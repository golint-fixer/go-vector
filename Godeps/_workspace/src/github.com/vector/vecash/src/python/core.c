#include <Python.h>
#include <alloca.h>
#include <stdint.h>
#include <stdlib.h>
#include <time.h>
#include "../libvecash/vecash.h"
#include "../libvecash/internal.h"

#if PY_MAJOR_VERSION >= 3
#define PY_STRING_FORMAT "y#"
#define PY_CONST_STRING_FORMAT "y"
#else
#define PY_STRING_FORMAT "s#"
#define PY_CONST_STRING_FORMAT "s"
#endif

#define MIX_WORDS (VECASH_MIX_BYTES/4)

static PyObject *
mkcache_bytes(PyObject *self, PyObject *args) {
    unsigned long block_number;
    unsigned long cache_size;

    if (!PyArg_ParseTuple(args, "k", &block_number))
        return 0;

    vecash_light_t L = vecash_light_new(block_number);
    PyObject * val = Py_BuildValue(PY_STRING_FORMAT, L->cache, L->cache_size);
    free(L->cache);
    return val;
}

/*
static PyObject *
calc_dataset_bytes(PyObject *self, PyObject *args) {
    char *cache_bytes;
    unsigned long full_size;
    int cache_size;

    if (!PyArg_ParseTuple(args, "k" PY_STRING_FORMAT, &full_size, &cache_bytes, &cache_size))
        return 0;

    if (full_size % MIX_WORDS != 0) {
        char error_message[1024];
        sprintf(error_message, "The size of data set must be a multiple of %i bytes (was %lu)", MIX_WORDS, full_size);
        PyErr_SetString(PyExc_ValueError, error_message);
        return 0;
    }

    if (cache_size % VECASH_HASH_BYTES != 0) {
        char error_message[1024];
        sprintf(error_message, "The size of the cache must be a multiple of %i bytes (was %i)", VECASH_HASH_BYTES, cache_size);
        PyErr_SetString(PyExc_ValueError, error_message);
        return 0;
    }

    vecash_params params;
    params.cache_size = (size_t) cache_size;
    params.full_size = (size_t) full_size;
    vecash_cache cache;
    cache.mem = (void *) cache_bytes;
    void *mem = malloc(params.full_size);
    vecash_compute_full_data(mem, &params, &cache);
    PyObject * val = Py_BuildValue(PY_STRING_FORMAT, (char *) mem, full_size);
    free(mem);
    return val;
}*/

// hashimoto_light(full_size, cache, header, nonce)
static PyObject *
hashimoto_light(PyObject *self, PyObject *args) {
    char *cache_bytes;
    char *header;
    unsigned long block_number;
    unsigned long long nonce;
    int cache_size, header_size;
    if (!PyArg_ParseTuple(args, "k" PY_STRING_FORMAT PY_STRING_FORMAT "K", &block_number, &cache_bytes, &cache_size, &header, &header_size, &nonce))
        return 0;
    if (header_size != 32) {
        char error_message[1024];
        sprintf(error_message, "Seed must be 32 bytes long (was %i)", header_size);
        PyErr_SetString(PyExc_ValueError, error_message);
        return 0;
    }
    struct vecash_light *s;
    s = calloc(sizeof(*s), 1);
    s->cache = cache_bytes;
    s->cache_size = cache_size;
    s->block_number = block_number;
    struct vecash_h256 *h;
    h = calloc(sizeof(*h), 1);
    for (int i = 0; i < 32; i++) h->b[i] = header[i];
    struct vecash_return_value out = vecash_light_compute(s, *h, nonce);
    return Py_BuildValue("{" PY_CONST_STRING_FORMAT ":" PY_STRING_FORMAT "," PY_CONST_STRING_FORMAT ":" PY_STRING_FORMAT "}",
                         "mix digest", &out.mix_hash, 32,
                         "result", &out.result, 32);
}
/*
// hashimoto_full(dataset, header, nonce)
static PyObject *
hashimoto_full(PyObject *self, PyObject *args) {
    char *full_bytes;
    char *header;
    unsigned long long nonce;
    int full_size, header_size;

    if (!PyArg_ParseTuple(args, PY_STRING_FORMAT PY_STRING_FORMAT "K", &full_bytes, &full_size, &header, &header_size, &nonce))
        return 0;

    if (full_size % MIX_WORDS != 0) {
        char error_message[1024];
        sprintf(error_message, "The size of data set must be a multiple of %i bytes (was %i)", MIX_WORDS, full_size);
        PyErr_SetString(PyExc_ValueError, error_message);
        return 0;
    }

    if (header_size != 32) {
        char error_message[1024];
        sprintf(error_message, "Header must be 32 bytes long (was %i)", header_size);
        PyErr_SetString(PyExc_ValueError, error_message);
        return 0;
    }


    vecash_return_value out;
    vecash_params params;
    params.full_size = (size_t) full_size;
    vecash_full(&out, (void *) full_bytes, &params, (vecash_h256_t *) header, nonce);
    return Py_BuildValue("{" PY_CONST_STRING_FORMAT ":" PY_STRING_FORMAT ", " PY_CONST_STRING_FORMAT ":" PY_STRING_FORMAT "}",
                         "mix digest", &out.mix_hash, 32,
                         "result", &out.result, 32);
}

// mine(dataset_bytes, header, difficulty_bytes)
static PyObject *
mine(PyObject *self, PyObject *args) {
    char *full_bytes;
    char *header;
    char *difficulty;
    srand(time(0));
    uint64_t nonce = ((uint64_t) rand()) << 32 | rand();
    int full_size, header_size, difficulty_size;

    if (!PyArg_ParseTuple(args, PY_STRING_FORMAT PY_STRING_FORMAT PY_STRING_FORMAT, &full_bytes, &full_size, &header, &header_size, &difficulty, &difficulty_size))
        return 0;

    if (full_size % MIX_WORDS != 0) {
        char error_message[1024];
        sprintf(error_message, "The size of data set must be a multiple of %i bytes (was %i)", MIX_WORDS, full_size);
        PyErr_SetString(PyExc_ValueError, error_message);
        return 0;
    }

    if (header_size != 32) {
        char error_message[1024];
        sprintf(error_message, "Header must be 32 bytes long (was %i)", header_size);
        PyErr_SetString(PyExc_ValueError, error_message);
        return 0;
    }

    if (difficulty_size != 32) {
        char error_message[1024];
        sprintf(error_message, "Difficulty must be an array of 32 bytes (only had %i)", difficulty_size);
        PyErr_SetString(PyExc_ValueError, error_message);
        return 0;
    }

    vecash_return_value out;
    vecash_params params;
    params.full_size = (size_t) full_size;

    // TODO: Multi threading?
    do {
        vecash_full(&out, (void *) full_bytes, &params, (const vecash_h256_t *) header, nonce++);
        // TODO: disagrees with the spec https://github.com/vector/wiki/wiki/Ethash#mining
    } while (!vecash_check_difficulty(&out.result, (const vecash_h256_t *) difficulty));

    return Py_BuildValue("{" PY_CONST_STRING_FORMAT ":" PY_STRING_FORMAT ", " PY_CONST_STRING_FORMAT ":" PY_STRING_FORMAT ", " PY_CONST_STRING_FORMAT ":K}",
            "mix digest", &out.mix_hash, 32,
            "result", &out.result, 32,
            "nonce", nonce);
}
*/

//get_seedhash(block_number)
static PyObject *
get_seedhash(PyObject *self, PyObject *args) {
    unsigned long block_number;
    if (!PyArg_ParseTuple(args, "k", &block_number))
        return 0;
    if (block_number >= VECASH_EPOCH_LENGTH * 2048) {
        char error_message[1024];
        sprintf(error_message, "Block number must be less than %i (was %lu)", VECASH_EPOCH_LENGTH * 2048, block_number);

        PyErr_SetString(PyExc_ValueError, error_message);
        return 0;
    }
    vecash_h256_t seedhash = vecash_get_seedhash(block_number);
    return Py_BuildValue(PY_STRING_FORMAT, (char *) &seedhash, 32);
}

static PyMethodDef PyvecashMethods[] =
        {
                {"get_seedhash", get_seedhash, MVEC_VARARGS,
                        "get_seedhash(block_number)\n\n"
                                "Gets the seedhash for a block."},
                {"mkcache_bytes", mkcache_bytes, MVEC_VARARGS,
                        "mkcache_bytes(block_number)\n\n"
                                "Makes a byte array for the cache for given block number\n"},
                /*{"calc_dataset_bytes", calc_dataset_bytes, MVEC_VARARGS,
                        "calc_dataset_bytes(full_size, cache_bytes)\n\n"
                                "Makes a byte array for the dataset for a given size given cache bytes"},*/
                {"hashimoto_light", hashimoto_light, MVEC_VARARGS,
                        "hashimoto_light(block_number, cache_bytes, header, nonce)\n\n"
                                "Runs the hashimoto hashing function just using cache bytes. Takes an int (full_size), byte array (cache_bytes), another byte array (header), and an int (nonce). Returns an object containing the mix digest, and hash result."},
                /*{"hashimoto_full", hashimoto_full, MVEC_VARARGS,
                        "hashimoto_full(dataset_bytes, header, nonce)\n\n"
                                "Runs the hashimoto hashing function using the dataset bytes. Useful for testing. Returns an object containing the mix digest (byte array), and hash result (another byte array)."},
                {"mine", mine, MVEC_VARARGS,
                        "mine(dataset_bytes, header, difficulty_bytes)\n\n"
                                "Mine for an adequate header. Returns an object containing the mix digest (byte array), hash result (another byte array) and nonce (an int)."},*/
                {NULL, NULL, 0, NULL}
        };

#if PY_MAJOR_VERSION >= 3
static struct PyModuleDef PyvecashModule = {
    PyModuleDef_HEAD_INIT,
    "pyvecash",
    "...",
    -1,
    PyvecashMethods
};

PyMODINIT_FUNC PyInit_pyvecash(void) {
    PyObject *module =  PyModule_Create(&PyvecashModule);
    // Following Spec: https://github.com/vector/wiki/wiki/Ethash#definitions
    PyModule_AddIntConstant(module, "REVISION", (long) VECASH_REVISION);
    PyModule_AddIntConstant(module, "DATASET_BYTES_INIT", (long) VECASH_DATASET_BYTES_INIT);
    PyModule_AddIntConstant(module, "DATASET_BYTES_GROWTH", (long) VECASH_DATASET_BYTES_GROWTH);
    PyModule_AddIntConstant(module, "CACHE_BYTES_INIT", (long) VECASH_CACHE_BYTES_INIT);
    PyModule_AddIntConstant(module, "CACHE_BYTES_GROWTH", (long) VECASH_CACHE_BYTES_GROWTH);
    PyModule_AddIntConstant(module, "EPOCH_LENGTH", (long) VECASH_EPOCH_LENGTH);
    PyModule_AddIntConstant(module, "MIX_BYTES", (long) VECASH_MIX_BYTES);
    PyModule_AddIntConstant(module, "HASH_BYTES", (long) VECASH_HASH_BYTES);
    PyModule_AddIntConstant(module, "DATASET_PARENTS", (long) VECASH_DATASET_PARENTS);
    PyModule_AddIntConstant(module, "CACHE_ROUNDS", (long) VECASH_CACHE_ROUNDS);
    PyModule_AddIntConstant(module, "ACCESSES", (long) VECASH_ACCESSES);
    return module;
}
#else
PyMODINIT_FUNC
initpyvecash(void) {
    PyObject *module = Py_InitModule("pyvecash", PyvecashMethods);
    // Following Spec: https://github.com/vector/wiki/wiki/Ethash#definitions
    PyModule_AddIntConstant(module, "REVISION", (long) VECASH_REVISION);
    PyModule_AddIntConstant(module, "DATASET_BYTES_INIT", (long) VECASH_DATASET_BYTES_INIT);
    PyModule_AddIntConstant(module, "DATASET_BYTES_GROWTH", (long) VECASH_DATASET_BYTES_GROWTH);
    PyModule_AddIntConstant(module, "CACHE_BYTES_INIT", (long) VECASH_CACHE_BYTES_INIT);
    PyModule_AddIntConstant(module, "CACHE_BYTES_GROWTH", (long) VECASH_CACHE_BYTES_GROWTH);
    PyModule_AddIntConstant(module, "EPOCH_LENGTH", (long) VECASH_EPOCH_LENGTH);
    PyModule_AddIntConstant(module, "MIX_BYTES", (long) VECASH_MIX_BYTES);
    PyModule_AddIntConstant(module, "HASH_BYTES", (long) VECASH_HASH_BYTES);
    PyModule_AddIntConstant(module, "DATASET_PARENTS", (long) VECASH_DATASET_PARENTS);
    PyModule_AddIntConstant(module, "CACHE_ROUNDS", (long) VECASH_CACHE_ROUNDS);
    PyModule_AddIntConstant(module, "ACCESSES", (long) VECASH_ACCESSES);
}
#endif
