// Copyright 2020 Redpanda Data, Inc.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.md
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0

#include "kafka/server/handlers/join_group.h"

#include "kafka/server/group_manager.h"
#include "kafka/server/group_router.h"
#include "kafka/server/request_context.h"
#include "kafka/server/response.h"
#include "utils/remote.h"
#include "utils/to_string.h"

#include <seastar/core/coroutine.hh>

namespace kafka {

template<>
process_result_stages join_group_handler::handle(
  request_context ctx, [[maybe_unused]] ss::smp_service_group g) {
    join_group_request request;
    request.decode(ctx.reader(), ctx.header().version);
    request.version = ctx.header().version;
    if (ctx.header().client_id) {
        request.client_id = kafka::client_id(
          ss::sstring(*ctx.header().client_id));
    }
    request.client_host = kafka::client_host(
      fmt::format("{}", ctx.connection()->client_host()));
    log_request(ctx.header(), request);

    if (!ctx.authorized(security::acl_operation::read, request.data.group_id)) {
        return process_result_stages::single_stage(ctx.respond(
          join_group_response(error_code::group_authorization_failed)));
    }

    auto stages = ctx.groups().join_group(std::move(request));
    auto res = ss::do_with(
      std::move(ctx),
      [f = std::move(stages.result)](request_context& ctx) mutable {
          return f.then([&ctx](join_group_response response) {
              return ctx.respond(std::move(response));
          });
      });

    return process_result_stages(std::move(stages.dispatched), std::move(res));
}

} // namespace kafka
