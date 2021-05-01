cmake_minimum_required(VERSION 3.10)

# set the project name
project(FilterSeed)
set(LIBRARIES gcrypt gpg-error cubiomes m pthread minecraft_nether_gen_rs)

message("Generator: ${CMAKE_GENERATOR}")
if(${CMAKE_SYSTEM_NAME} MATCHES "Windows.*")
    set(LIB_DIR libs/windows_x86_64)
    set(EXECUTABLE_POSTFIX ".exe")
    list(APPEND LIBRARIES userenv intl.dll ws2_32)
else()
    set(LIB_DIR libs/linux_x86_64)
    set(EXECUTABLE_POSTFIX "")
    list(APPEND LIBRARIES dl)
endif()


link_libraries(${LIBRARIES})
include_directories(include)
link_directories(${LIB_DIR})

file(GLOB_RECURSE SOURCE_FILES */csprng.c)
file(MAKE_DIRECTORY ${CMAKE_BINARY_DIR}/generated)
foreach(SOURCE_FILE IN LISTS SOURCE_FILES)
    string(REGEX REPLACE ".*/([^/]+)/csprng.c" "\\1" FOLDER_NAME ${SOURCE_FILE})
    set(PATCHED_SOURCE ${CMAKE_BINARY_DIR}/generated/${FOLDER_NAME}.c)
    add_custom_command(
        OUTPUT ${PATCHED_SOURCE}
        COMMAND sed -r "s/unsigned long/uint64_t/g; s/([0-9a-fA-F])UL\\b/\\1ULL/g; s/long/int64_t/g; s/Seed: %ld/Seed: %lld/g;" ${SOURCE_FILE} > ${PATCHED_SOURCE}
	    COMMAND patch ${PATCHED_SOURCE} ${CMAKE_SOURCE_DIR}/rand.patch
        DEPENDS ${SOURCE_FILE}
        VERBATIM)
    message("Find build target ${FOLDER_NAME}")
    add_executable(${FOLDER_NAME} ${PATCHED_SOURCE})
    set(INSTALL_DIR ${CMAKE_BINARY_DIR}/fsg/${FOLDER_NAME})
    file(MAKE_DIRECTORY ${INSTALL_DIR})
    install(
        FILES ${CMAKE_BINARY_DIR}/${FOLDER_NAME}${EXECUTABLE_POSTFIX}
        DESTINATION ${INSTALL_DIR}
        RENAME seed${EXECUTABLE_POSTFIX}
        PERMISSIONS OWNER_EXECUTE
    )
    install(
        FILES ${SOURCE_FILE}
        DESTINATION ${INSTALL_DIR}
    )
    install(
        FILES runner/run.bat
        DESTINATION ${INSTALL_DIR}
        PERMISSIONS OWNER_EXECUTE        
    )
endforeach()
